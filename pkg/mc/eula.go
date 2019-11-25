package mc

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EULA indicates the user acceptance of the Minecraft EULA.
type EULA struct {
	*viper.Viper
	EULA bool
}

type eulaData struct {
	timestamp time.Time
}

const eulaTemplate = `#By changing the setting below to TRUE you are indicating your agreement to our EULA (https://account.mojang.com/documents/minecraft_eula).
#{{ .Timestamp }}
eula=false
`

// CheckEULA checks for the presence of a file named eula.txt in the current
// directory, containing the line "eula=true".
func CheckEULA() error {
	eula := &EULA{}
	eula.Viper = viper.New()

	eula.SetConfigFile("eula.txt")
	eula.SetConfigType("properties")

	if err := eula.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return err
		}

		file, cerr := os.OpenFile(eula.ConfigFileUsed(), os.O_RDWR|os.O_CREATE, 0644)
		if cerr != nil {
			return cerr
		}

		if werr := writeEULA(file); werr != nil {
			return werr
		}
	}

	if err := eula.Unmarshal(eula); err != nil {
		return err
	}

	if !eula.EULA {
		return fmt.Errorf("You need to agree to the EULA in order to run the server. Go to eula.txt for more info.")
	}

	return nil
}

func newEULAData() *eulaData {
	return &eulaData{
		timestamp: time.Now(),
	}
}

func (data *eulaData) Timestamp() string {
	return data.timestamp.Format("Mon Jan 02 15:04:05 MST 2006")
}

func writeEULA(wr io.Writer) error {
	log.Debug("Writing EULA file")

	tmpl, err := template.New("eula").Parse(eulaTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, newEULAData()); err != nil {
		return err
	}

	return nil
}
