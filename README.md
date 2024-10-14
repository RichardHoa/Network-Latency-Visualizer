# Network-Latency-Visualizer

- Network Latency Checks: The program will regularly check network latency using the ping command. Users can set the frequency for these checks, ranging from every 5 minutes to once a day, utilizing a cron job to automate this process.
- Data Visualization: After gathering latency data over time, users will have the option to visualize this information in a chart format.
- Performance Recommendations: The application will provide performance recommendations, such as identifying which hours have the strongest network performance.
- Download and Upload Speed: Users will also be able to view their download and upload speed.
- Bandwidth usage for each process (program)


We use ping to check for network latency, please note that this only measure the latency of the whole round trip, meaning we can't know whether it's the upload part that is slow or the download part that is slow in your machine. We ping 10 times to check for the result


1. chmod 777 scanning
2. ./scanning