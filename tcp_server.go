package main

import (
	"io"
	"log"
	"net"
	"sync"
	"strings"
	"strconv"
	"fmt"
	"errors"
	"github.com/xBroccoliMaster69x/funtemps/conv"
	"github.com/xBroccoliMaster69x/is105sem03/mycrypt"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
  				        case "ping":
						kryptertMelding := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03,4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))
					default:
						if strings.HasPrefix(string(dekryptertMelding), "Kjevik"){
							fahrenheitLine, _ := CelsiusToFahrenheitLine(string(dekryptertMelding))
							kryptertMelding := mycrypt.Krypter([]rune(fahrenheitLine), mycrypt.ALF_SEM03, 4)
							_, err = c.Write([]byte(string(kryptertMelding)))
						} else {
							kryptertMelding := mycrypt.Krypter([]rune(dekryptertMelding), mycrypt.ALF_SEM03, 4)
							log.Println("kryptert melding: ", string(kryptertMelding))
							_, err = c.Write([]byte(string(kryptertMelding)))
						}
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
func CelsiusToFahrenheitString(celsius string) (string, error) {
	var fahrFloat float64
	var err error
	if celsiusFloat, err := strconv.ParseFloat(celsius, 64); err == nil {
		fahrFloat = conv.CelsiusToFahrenheit(celsiusFloat)
	}
	fahrString := fmt.Sprintf("%.1f", fahrFloat)
	return fahrString, err
}

func CelsiusToFahrenheitLine(line string) (string, error) {
	dividedString := strings.Split(line, ";")
	var err error

	if len(dividedString) == 4 {
		if dividedString[3] != "" {

			dividedString[3], err = CelsiusToFahrenheitString(dividedString[3])
			if err != nil {
				return "", err
			}
		} else {
			return "Data er basert paa gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Alexander Glasdam Andersen;;", nil
		}
	} else {
		return "", errors.New("linje har ikke forventet format")
	}
	return strings.Join(dividedString, ";"), nil
}
