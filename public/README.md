# Setting up

You need:
- http-server (or an alternative HTTP server like nginx to serve the webpage)

Create a folder named 'labels' in the 'public' directory and add photos of every individual you want
to recognize. The name of the photo also becomes the label for the individual.

In 'workshop.js' on line 35, you'll find this line:

```
const labels = ['sjaak', 'richard', 'arvid', 'mara', 'frank', 'wilma', 'marcin', 'deha', 'slave']
```

Remove the names in the array and replace them with the names that you've provided photos for in your
'labels' directory. Format being '<name>.png'.

Open the website with http-server (or any alternative), allow the webcam feed to start and the algorithm
should start recognising the individuals that you've provided labels for.
