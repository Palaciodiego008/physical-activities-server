package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Activity struct {
	Id                            int     `json:"id"`
	UserId                        int     `json:"userId"`
	StartTimeInSeconds            int     `json:"startTimeInSeconds"`
	DurationInSeconds             int     `json:"durationInSeconds"`
	DistanceInMeters              float64 `json:"distanceInMeters"`
	Steps                         int     `json:"steps"`
	AverageSpeedInMetersPerSecond float64 `json:"averageSpeedInMetersPerSecond"`
	AveragePaceInMinutesPerKm     float64 `json:"averagePaceInMinutesPerKm"`
	TotalElevationGainInMeters    float64 `json:"totalElevationGainInMeters"`
	AverageHeartRateInBPM         int     `json:"averageHeartRateInBPM"`
}

func readActivities(filename string) ([]Activity, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var activities []Activity
	activityCh := make(chan Activity)

	// Use WaitGroup para esperar que todas las goroutines terminen
	var wg sync.WaitGroup

	// Función para procesar una línea y enviar el resultado a los canales
	processLine := func(line []string) {
		defer wg.Done()

		id, _ := strconv.Atoi(line[0])
		userId, _ := strconv.Atoi(line[1])
		startTimeInSeconds, _ := strconv.Atoi(line[2])
		durationInSeconds, _ := strconv.Atoi(line[3])
		distanceInMeters, _ := strconv.ParseFloat(line[4], 64)
		steps, _ := strconv.Atoi(line[5])
		averageSpeed, _ := strconv.ParseFloat(line[6], 64)
		averagePace, _ := strconv.ParseFloat(line[7], 64)
		elevationGain, _ := strconv.ParseFloat(line[8], 64)
		heartRate, _ := strconv.Atoi(line[9])

		activity := Activity{
			Id:                            id,
			UserId:                        userId,
			StartTimeInSeconds:            startTimeInSeconds,
			DurationInSeconds:             durationInSeconds,
			DistanceInMeters:              distanceInMeters,
			Steps:                         steps,
			AverageSpeedInMetersPerSecond: averageSpeed,
			AveragePaceInMinutesPerKm:     averagePace,
			TotalElevationGainInMeters:    elevationGain,
			AverageHeartRateInBPM:         heartRate,
		}

		// Envía la actividad a través del canal
		activityCh <- activity
	}

	// Use goroutines para procesar líneas en paralelo
	for _, line := range lines {
		wg.Add(1)
		go processLine(line)
	}

	// Cierre el canal de actividades después de que todas las goroutines hayan terminado
	go func() {
		wg.Wait()
		close(activityCh)
	}()

	// Recoge actividades del canal y agrégales al slice
	for activity := range activityCh {
		activities = append(activities, activity)
	}

	return activities, nil
}

func detectSuspiciousActivities(activities []Activity) []Activity {
	var suspiciousActivities []Activity

	// Calcula estadísticas para comparación
	var (
		averagePace      float64
		averageDuration  float64
		averageDistance  float64
		maxElevationGain float64
	)

	for _, activity := range activities {
		averagePace += activity.AveragePaceInMinutesPerKm
		averageDuration += float64(activity.DurationInSeconds)
		averageDistance += activity.DistanceInMeters
		if activity.TotalElevationGainInMeters > maxElevationGain {
			maxElevationGain = activity.TotalElevationGainInMeters
		}
	}

	numActivities := float64(len(activities))
	averagePace /= numActivities
	averageDuration /= numActivities
	averageDistance /= numActivities

	// Establece umbrales para detección
	paceThreshold := 4.0                  // ejemplo, ajusta según sea necesario
	elevationThreshold := 2.0             // ejemplo, ajusta según sea necesario
	durationDistanceRatioThreshold := 2.0 // ejemplo, ajusta según sea necesario

	for _, activity := range activities {
		// Detecta actividades sospechosas
		if activity.AveragePaceInMinutesPerKm < paceThreshold ||
			activity.TotalElevationGainInMeters > elevationThreshold*maxElevationGain ||
			(float64(activity.DurationInSeconds)/averageDuration) > durationDistanceRatioThreshold ||
			(activity.DistanceInMeters/averageDistance) > durationDistanceRatioThreshold {
			suspiciousActivities = append(suspiciousActivities, activity)
		}
	}

	return suspiciousActivities
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error al obtener el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crear un directorio para almacenar archivos cargados
	uploadDir := "./uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModeDir)
	}

	// Crear un nombre de archivo único
	fileName := uploadDir + strconv.FormatInt(time.Now().Unix(), 10) + "_" + handler.Filename

	// Guardar el archivo en el servidor
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	json.NewEncoder(w).Encode(map[string]string{"fileName": fileName})
}

// Modifica la función getActivitiesHandler para leer el archivo más reciente
func getActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener la lista de archivos en el directorio de carga
	dirEntries, err := os.ReadDir("./uploads")
	if err != nil {
		http.Error(w, "Error al obtener la lista de archivos", http.StatusInternalServerError)
		return
	}

	// Estructura para almacenar información de archivos y su índice correspondiente
	type fileInfo struct {
		index int
		info  os.FileInfo
	}

	// Crear una slice de fileInfo
	var filesWithInfo []fileInfo

	// Obtener información de cada archivo
	for i, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			http.Error(w, "Error al obtener información del archivo", http.StatusInternalServerError)
			return
		}
		filesWithInfo = append(filesWithInfo, fileInfo{index: i, info: info})
	}

	// Ordenar los archivos por fecha de modificación (el más reciente primero)
	sort.Slice(filesWithInfo, func(i, j int) bool {
		return filesWithInfo[i].info.ModTime().After(filesWithInfo[j].info.ModTime())
	})

	// Tomar el archivo más reciente
	if len(filesWithInfo) == 0 {
		http.Error(w, "No hay archivos disponibles", http.StatusNotFound)
		return
	}
	fileName := "./uploads/" + dirEntries[filesWithInfo[0].index].Name()

	// Leer las actividades desde el archivo más reciente
	activities, err := readActivities(fileName)
	if err != nil {
		http.Error(w, "Error al leer las actividades", http.StatusInternalServerError)
		return
	}

	// Detectar actividades sospechosas
	suspiciousActivities := detectSuspiciousActivities(activities)

	// Devolver las actividades (sospechosas o no) como respuesta JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suspiciousActivities); err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Configura el enrutador Gorilla Mux
	r := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Endoints
	r.HandleFunc("/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/activities", getActivitiesHandler).Methods("GET")

	// Aplica el middleware CORS a todas las rutas
	corsRouter := handlers.CORS(headers, methods, origins)(r)

	// Inicia el servidor en el puerto 8080
	port := 8080
	fmt.Printf("Server listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), corsRouter))
}
