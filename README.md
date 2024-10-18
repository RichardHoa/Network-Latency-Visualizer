# Network Latency Visualizer

## Features
1. **Data Visualization**: Visualize network latency data over time in easy-to-understand charts.
2. **Network Latency Monitoring**: Perform regular network latency checks using the `ping` command. You can set the frequency of checks (from every 5 minutes to once a day) via a cron job to automate the process.
3. **Process Bandwidth Usage**: Track bandwidth usage by individual processes, displaying both incoming and outgoing data.
4. **Download & Upload Speed**: Measure and display your current download and upload speeds.

## How to Use

1. Make the script executable:
    ```bash
    chmod 777 scanning
    ```

2. Start data collection:
    ```bash
    ./scanning
    ```

3. For advanced options, run:
    ```bash
    ./scanning -a
    ```

### Advanced Options Menu
Upon running the advanced options, you'll see the following menu:

A terminal option screen will appear with various option, pick what you like!
```
What do you want to do?:
>  Cronjob options
   Show network bandwidth consumed by top 3 processes
   Show network latency chart
   Speed testing
   Quit
```

#### Options Explained
- **Cronjob Options**: Modify or remove the existing cronjob that automates network checks. You can view the current cronjob by using:
    ```bash
    crontab -l
    ```
    The cronjob working directory will resemble: `$Yourworkingdir/go-networking/scanning`.
  
- **Show Network Bandwidth by Top 3 Processes**: Displays two HTML graphs in your browser:
    1. Incoming network bandwidth usage by the top 3 processes.
    2. Outgoing network bandwidth usage by the top 3 processes.

   Additionally, a detailed table with all processes and their bandwidth usage will be displayed in the terminal.

- **Show Network Latency Chart**: Opens an HTML graph with a full network latency performance overview.

All HTML charts are stored in the `chart/html` folder for future access.

### Data Storage
- **Network Bandwidth Data**: Stored in `network/network.txt`.
- **Network Latency Data**: Stored in `ping/ping.txt`.

## Motives
This project was created as a way to get familiar with the Go programming language, combined with an interest in networking.

## Technology
- **Latency Data Collection**: Uses the built-in macOS `ping` command (`ping google.com -c 10`) to gather latency data.
- **Bandwidth Usage**: Uses `nettop -l 1 -P -x` to monitor bandwidth usage by each process.

## Limitations
- **Latency Measurement**: The `ping` command only measures the total round-trip latency, so it cannot distinguish whether upload or download is slower.
- **Process Name Length**: The `nettop` command truncates long process names, but it's usually clear enough to identify the associated application.
