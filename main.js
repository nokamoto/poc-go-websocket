
const socket = new WebSocket('ws://localhost:9090/ws');

socket.addEventListener('open', function (event) {
    socket.send('Hello Server!');
});

var number = 0;

window.onload = function() {
    navigator.mediaDevices.getUserMedia({ video: false, audio: true })
    .then(stream => {
        console.log(stream);
        // var video = document.getElementById('v');
        // video.srcObject = stream;

        var mediaRecorder = new MediaRecorder(stream);
        mediaRecorder.start(1000);
        mediaRecorder.ondataavailable = function(e) {
            console.log(e);
            console.log(e.data);
            socket.send(e.data);
        }
    })
    .catch(err => {
        console.log(err);
    });
};
