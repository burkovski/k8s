package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
)

func storeData(ctx echo.Context) error {
	file, err := os.Create(dataFilename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Unable to create data file: %s", err))
	}
	defer file.Close()

	if _, err := io.Copy(file, ctx.Request().Body); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Unable to copy data from request to data file: %s", err))
	}

	log.Info().Str("filename", dataFilename).
		Msg("new data has been received and stored")

	hostname, err := os.Hostname()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Unable to get hostname: %s", err))
	}

	response := fmt.Sprintf("Data stored pod %q", hostname)

	return ctx.String(http.StatusCreated, response)
}

func getData(ctx echo.Context) error {
	dataFromFile, err := os.ReadFile(dataFilename)
	if os.IsNotExist(err) {
		dataFromFile = []byte("No data posted yet")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Unable to read data from data file: %s", err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Unable to get hostname: %s", err))
	}

	response := fmt.Sprintf("You've hit: %q\nData stored on this pod: %q",
		hostname, string(dataFromFile))
	return ctx.String(http.StatusOK, response)
}
