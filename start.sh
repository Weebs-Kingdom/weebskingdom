cd ./main
echo building server...
/usr/local/go/bin/go build
echo server built!
echo starting server...
screen -AdmS newweeb ./weebskingdom
echo server started!
