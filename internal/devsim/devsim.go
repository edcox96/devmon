package devsim

import (
    "errors"
	"log"
)

/* easier to use devsim as a package devsim rather than a stanalone command
func main() {
    log.Printf("devsim main\n")
    if err := InitDevSim(); err != nil {
        log.Fatalf("InitDevSim failed! %s", err)
    }
}
*/

func InitDevSim() error {
    sim := NewUsbDevSim()
    con, err := sim.NewConsole("Mercy Hospital",
                          "0001-0000",
                          1)
    if err != nil {
        log.Printf("sim.NewConsole failed! %s", err)
        return err
    }

    sim.Consoles = append(sim.Consoles, con)

    return nil
}

func NewUsbDevSim() *UsbDevSim {
	sim := &usbDevSim
    sim.Consoles = make([]*Console, 0, 2)
	return sim
}

var usbDevSim = UsbDevSim {}

type UsbDevSim struct {
    Consoles []*Console
}

func (sim *UsbDevSim) NewConsole(name string, sn string, id uint64) (*Console, error) {
    var con = new(Console)
    con.Name = name
    con.SN = sn
    con.ID = id
    return con, nil
}

type Console struct {
    Name string
    SN   string
    ID   uint64
}

func GetConsole(id uint64) (*Console, error) {
    sim := &usbDevSim
    if int(id) > len(sim.Consoles) {
        log.Printf("id %d is invalid!", id)
        return nil, errors.New("id is invalid")
    }
    con := sim.Consoles[id - 1]
    return con, nil
}