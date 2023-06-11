# myip

Small utility to find one's IP. Can also be used a as library!
Get it from the github releases or build it yourself with

```bash
go clone https://github.com/wnke/myip
cd myip
go build -o myip ./cmd
cp myip $HOME/.local/bin/myip
```

## Usage

### Command line

```bash
$ myip
123.123.111.111
```

### Library

```golang
import 	"github.com/wnke/myip"

//...
ipa, err := myip.NewIPDiscover()
if err != nil {
    return err
}

ip, err := ipa.Discover()
if err != nil {
    return err
}
fmt.Print(ip.String())

```
