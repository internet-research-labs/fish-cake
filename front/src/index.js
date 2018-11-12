import {PlanetApp} from './PlanetApp.js';


/**
 * Game State
 */
let STATE = {
  ships: [],
};


function lezgo(id) {

  let app = new PlanetApp({
    id: id,
    state: STATE,
  });

  app.setup();

  // Get live stream
  let socket = new WebSocket("ws://localhost:8000/ws");

  fetch("/world").then((resp) => {
    return resp.json();
  }).then((world) => {
    app.updateWorld(world);
  });

  socket.onopen = function (a, b, c) {
    console.log("OPENED!", a, b, c);

    socket.addEventListener('message', function (ev) {
      let data = JSON.parse(ev.data);
      let ships = data.ships;
      app.updateShips(ships);
      app.needsUpdate = true;
    });

    socket.addEventListener("close", function () {
      console.log("CLOSED!");
    });
  }

  // Render loop
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
