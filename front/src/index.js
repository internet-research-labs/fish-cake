import {PlanetApp} from './PlanetApp.js';


function lezgo(id) {

  // Initial game state
  fetch("/world").then((resp) => {
    return resp.json();
  }).then((world) => {
    console.log(world);
  });

  // Get live stream
  let socket = new WebSocket("ws://localhost:8000/ws");

  socket.onopen = function (a, b, c) {
    console.log("OPENED!", a, b, c);

    socket.addEventListener('message', function (ev) {
      console.log("MESSAGE:", ev);
    });

    socket.addEventListener("close", function () {
      console.log("CLOSED!");
    });
  }

  // Render loop
  let app = new PlanetApp({
    id: id,
  });

  app.setup();

  (function loop() {
    app.update();
    app.draw();
    requestAnimationFrame(loop);
  }());


  // Things that you might want later
  return [
    app,
    socket,
  ];
}

export {
  lezgo,
}
