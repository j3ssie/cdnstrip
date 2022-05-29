## CDNStrip

Striping CDN IPs from a list of IP Addresses

Built based on the [projectdiscovery/cdncheck](github.com/projectdiscovery/cdncheck) library.

## Install

```shell
go install github.com/j3ssie/cdnstrip@latest
```

## Usage

```shell
# simple usage
cat ips | cdnstrip -c 50

# write the output to a file
cat ips | cdnstrip -cdn cdn.txt -n non-cdn.txt
```

## Donation!!!

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://paypal.me/j3ssiejjj)

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/j3ssie)
