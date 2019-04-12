
const socket = new WebSocket('ws://localhost:9090/ws');

socket.addEventListener('open', function (event) {
    socket.send('Hello Server!');
});

socket.addEventListener('message', function (event) {
    console.log('Message from server ', event.data);
});

var number = 0;

setInterval(function() {
    number = number + 1;
    socket.send('echo #' + number);
}, 3000);
