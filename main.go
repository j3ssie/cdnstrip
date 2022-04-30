package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/projectdiscovery/cdncheck"
    "log"
    "net"
    "net/url"
    "os"
    "strings"
    "sync"
)

// cat /tmp/list_of_IP | cdnstrip -c 100
var (
    concurrency int
    verbose     bool
    writeOutput bool
    nonCdnOut   string
    cdnOut      string
)

var cdnClient *cdncheck.Client
var nonCdnOutputWriter *os.File
var cdnOutputWriter *os.File

func main() {
    // cli arguments
    flag.IntVar(&concurrency, "c", 20, "Set the concurrency level")
    flag.StringVar(&nonCdnOut, "n", "", "Write non-CDN IPs to file")
    flag.StringVar(&cdnOut, "cdn", "", "Write CDN IPs to file")
    flag.BoolVar(&verbose, "v", false, "Verbose output with vendor of CDN")
    flag.Parse()

    var err error
    if cdnClient, err = cdncheck.NewWithCache(); err != nil {
        log.Fatal(err)
    }

    if nonCdnOut != "" {
        nonCdnOutputWriter, err = os.OpenFile(nonCdnOut, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to create/open noneCdnOutputFile\n")
            os.Exit(1)
        }
        defer nonCdnOutputWriter.Close()
        writeOutput = true
    }

    if cdnOut != "" {
        cdnOutputWriter, err = os.OpenFile(cdnOut, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to create/open notCdnOutputFile\n")
            os.Exit(1)
        }
        defer cdnOutputWriter.Close()
        writeOutput = true

    }


    var wg sync.WaitGroup
    jobs := make(chan string, concurrency)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                // actually start checking
                cdnChecking(job)
            }
        }()
    }

    sc := bufio.NewScanner(os.Stdin)
    go func() {
        for sc.Scan() {
            line := strings.TrimSpace(sc.Text())
            if err := sc.Err(); err == nil && line != "" {
                jobs <- line
            }
        }
        close(jobs)
    }()
    wg.Wait()
}

func cdnChecking(ip string) {
    // in case input as http format
    if strings.HasPrefix(ip, "http") {
        // parse url
        uu , err := url.Parse(ip)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to parse url: %s\n", err)
            return
        }
        ip = uu.Hostname()
    }

    if found, vendor, ok := cdnClient.Check(net.ParseIP(ip)); found && ok == nil {
        if writeOutput {
            nonCdnOutputWriter.WriteString(ip + "\n")
        }
    } else {
        line := ip
        if verbose {
            fmt.Println(vendor, ",", ip)
            line = fmt.Sprintf("%s,%s\n", vendor, ip)
        }

        fmt.Println(line)
        if writeOutput {
            cdnOutputWriter.WriteString(line + "\n")
        }
    }
}
