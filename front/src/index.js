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

  fetch("/world").then((resp) => {
    return resp.json();
  }).then((world) => {
    app.updateWorld(world);
    attachSocketEvents();
  });


  function attachSocketEvents() {

    // Get live stream
    let socket = new WebSocket("ws://"+location.host+"/ws");

    // Open
    socket.addEventListener('open', function (ev) {
      console.log(ev);
    });

    // Message
    socket.addEventListener('message', function (ev) {
      let data = JSON.parse(ev.data);
      let ships = data.ships;

      switch (data.Type) {
      case "you-are":
          console.log("You are ", data.Id);
          app.targetShip = data.Id;
      default:
        app.updateShips(ships);
        app.needsUpdate = true;
      }
    });

    // Close
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
  ];
}

export {
  lezgo,
}
