package model

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mounttypes "docker.io/go-docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	units "github.com/docker/go-units"
)

type Course struct {
	ID        int            `json:"id,omitempty"`
	Title     string         `json:"title,omitempty"`
	Comments  string         `json:"comments,omitempty"`
	Scenarios []ScenarioBase `json:"scenarios,omitempty"`
	Tags      []string       `json:"tags,omitempty"`
	Active    bool           `json:"active,omitempty"`
	Created   time.Time      `json:"created,omitempty"`
}
type ScenarioBase struct {
	Title    string `json:"title,omitempty"`
	Comments string `json:"comments,omitempty"`
	URL      string `json:"url,omitempty"`
}
type Scenario struct {
	// Title    string `json:"title,omitempty"`
	// Comments string `json:"comments,omitempty"`
	// URL      string `json:"url,omitempty"`
	ScenarioBase
	Created Time `json:"created,omitempty"`
	Intro   MarkDown
	Steps   []MarkDown
	Backend []Backend `json:"backend,omitempty"`
}

type Resource struct {
	ID       int    `json:"id,omitempty"`
	Course   string `json:"course,omitempty"`
	Scenario string `json:"scenario,omitempty"`
	Active   bool   `json:"active"`
	Updated  Time   `json:"updated,omitempty"`
}

type Backend struct {
	Command     string   `json:"cmd,omitempty"`
	Environment []string `json:"env,omitempty"`
	Image       string   `json:"image,omitempty"`
	Mounts      []string `json:"mounts,omitempty"`
	Expose      []string `json:"expose,omitempty"`
	PrivateIP   []string `json:"privateIP,omitempty"`
	Instances   int      `json:"instances,omitempty"`
}
type MarkDown struct {
	Title string `json:"title,omitempty"`
	File  string `json:"file,omitempty"`
}

// func (t Time) String() string {
//     return time.Time(t).Format(timeFormart)
// }

type MountOpt struct {
	values []mounttypes.Mount
}

// Set a new mount value
// nolint: gocyclo
func (m *MountOpt) Set(value string) error {
	csvReader := csv.NewReader(strings.NewReader(value))
	fields, err := csvReader.Read()
	if err != nil {
		return err
	}

	mount := mounttypes.Mount{}

	volumeOptions := func() *mounttypes.VolumeOptions {
		if mount.VolumeOptions == nil {
			mount.VolumeOptions = &mounttypes.VolumeOptions{
				Labels: make(map[string]string),
			}
		}
		if mount.VolumeOptions.DriverConfig == nil {
			mount.VolumeOptions.DriverConfig = &mounttypes.Driver{}
		}
		return mount.VolumeOptions
	}

	bindOptions := func() *mounttypes.BindOptions {
		if mount.BindOptions == nil {
			mount.BindOptions = new(mounttypes.BindOptions)
		}
		return mount.BindOptions
	}

	tmpfsOptions := func() *mounttypes.TmpfsOptions {
		if mount.TmpfsOptions == nil {
			mount.TmpfsOptions = new(mounttypes.TmpfsOptions)
		}
		return mount.TmpfsOptions
	}

	setValueOnMap := func(target map[string]string, value string) {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) == 1 {
			target[value] = ""
		} else {
			target[parts[0]] = parts[1]
		}
	}

	mount.Type = mounttypes.TypeVolume // default to volume mounts
	// Set writable as the default
	for _, field := range fields {
		parts := strings.SplitN(field, "=", 2)
		key := strings.ToLower(parts[0])

		if len(parts) == 1 {
			switch key {
			case "readonly", "ro":
				mount.ReadOnly = true
				continue
			case "volume-nocopy":
				volumeOptions().NoCopy = true
				continue
			}
		}

		if len(parts) != 2 {
			return fmt.Errorf("invalid field '%s' must be a key=value pair", field)
		}

		value := parts[1]
		switch key {
		case "type":
			mount.Type = mounttypes.Type(strings.ToLower(value))
		case "source", "src":
			mount.Source = value
		case "target", "dst", "destination":
			mount.Target = value
		case "readonly", "ro":
			mount.ReadOnly, err = strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, value)
			}
		case "consistency":
			mount.Consistency = mounttypes.Consistency(strings.ToLower(value))
		case "bind-propagation":
			bindOptions().Propagation = mounttypes.Propagation(strings.ToLower(value))
		case "volume-nocopy":
			volumeOptions().NoCopy, err = strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for volume-nocopy: %s", value)
			}
		case "volume-label":
			setValueOnMap(volumeOptions().Labels, value)
		case "volume-driver":
			volumeOptions().DriverConfig.Name = value
		case "volume-opt":
			if volumeOptions().DriverConfig.Options == nil {
				volumeOptions().DriverConfig.Options = make(map[string]string)
			}
			setValueOnMap(volumeOptions().DriverConfig.Options, value)
		case "tmpfs-size":
			sizeBytes, err := units.RAMInBytes(value)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, value)
			}
			tmpfsOptions().SizeBytes = sizeBytes
		case "tmpfs-mode":
			ui64, err := strconv.ParseUint(value, 8, 32)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %s", key, value)
			}
			tmpfsOptions().Mode = os.FileMode(ui64)
		default:
			return fmt.Errorf("unexpected key '%s' in '%s'", key, field)
		}
	}

	if mount.Type == "" {
		return fmt.Errorf("type is required")
	}

	if mount.Target == "" {
		return fmt.Errorf("target is required")
	}

	if mount.VolumeOptions != nil && mount.Type != mounttypes.TypeVolume {
		return fmt.Errorf("cannot mix 'volume-*' options with mount type '%s'", mount.Type)
	}
	if mount.BindOptions != nil && mount.Type != mounttypes.TypeBind {
		return fmt.Errorf("cannot mix 'bind-*' options with mount type '%s'", mount.Type)
	}
	if mount.TmpfsOptions != nil && mount.Type != mounttypes.TypeTmpfs {
		return fmt.Errorf("cannot mix 'tmpfs-*' options with mount type '%s'", mount.Type)
	}

	m.values = append(m.values, mount)
	return nil
}

// Type returns the type of this option
func (m *MountOpt) Type() string {
	return "mount"
}

// String returns a string repr of this option
func (m *MountOpt) String() string {
	mounts := []string{}
	for _, mount := range m.values {
		repr := fmt.Sprintf("%s %s %s", mount.Type, mount.Source, mount.Target)
		mounts = append(mounts, repr)
	}
	return strings.Join(mounts, ", ")
}

// Value returns the mounts
func (m *MountOpt) Value() []mounttypes.Mount {
	return m.values
}

func (s Scenario) Base() ScenarioBase {
	return ScenarioBase{Title: s.Title, URL: s.URL, Comments: s.Comments}
}

func (backend *Backend) GetMounts(template map[string]string) []mounttypes.Mount {
	var mountRef MountOpt

	for _, mountValue := range backend.Mounts {

		for k, v := range template {
			mountValue = strings.Replace(mountValue, "{"+k+"}", v, -1)

		}
		mountRef.Set(mountValue)

	}

	return mountRef.Value()
}

func ParsePortBidings(backend Backend) nat.PortMap {
	portMap := make(nat.PortMap)
	for _, v := range backend.Expose {
		portMap[nat.Port(v)] = []nat.PortBinding{}
		if len(backend.PrivateIP) == 0 {
			portMap[nat.Port(v)] = append(portMap[nat.Port(v)], nat.PortBinding{HostIP: "127.0.0.1", HostPort: ""})
		} else {
			for _, ip := range backend.PrivateIP {
				portMap[nat.Port(v)] = append(portMap[nat.Port(v)], nat.PortBinding{HostIP: ip, HostPort: ""})
			}
		}

	}
	log.Printf("%+v", portMap)

	return portMap
}

func ParseExposedPorts(ports nat.PortMap) nat.PortSet {
	exports := make(nat.PortSet, 10)
	for k, _ := range ports {
		exports[k] = struct{}{}
	}
	return exports
}
