import {PlanetApp} from './PlanetApp.js';
import {SwarmApp} from './SwarmApp.js';


/**
 * Game State
 */
let STATE = {
  ships: [],
};

let socket = undefined;

function svvarm(id) {
  let app = new SwarmApp({
    id: id,
    state: STATE,
  });

  app.setup();
  app.buildWorld();
  app.update();
  app.draw();

  // Get live stream
  socket = new WebSocket("ws://"+location.host+"/swarm");

  // Open
  socket.addEventListener('open', function (ev) {
    console.log(ev);
  });

  // Message
  socket.addEventListener('message', function (ev) {
    let data = JSON.parse(ev.data);
    let ships = data.ships;
    let type = data.type || data.Type;

    // Switch on this
    switch (type) {
      case "yupdate":
        app.updateSwifts(data.blob);
        app.update();
        app.draw();
    }
  });

  // Close
  socket.addEventListener("close", function () {
    console.log("CLOSED!");
  });
}

function updateSwarm(attraction, repulsion, alignment) {
  socket.send(JSON.stringify({
    "attraction": parseFloat(attraction),
    "repulsion": parseFloat(repulsion),
    "alignment": parseFloat(alignment),
  }));
}

function lezgo(id) {

  let app = new PlanetApp({
    id: id,
    state: STATE,
  });

  let {keypress} = app.eventHandlers();

  document.addEventListener("keypress", keypress);

  app.setup();

  attachSocketEvents();

  function attachSocketEvents() {

    let targetShip = undefined;

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
      let type = data.type || data.Type;

      // Switch on this
      switch (type) {
        case "world":
          app.updateWorld(data.World);
          break;

        case "you-are":
          app.id = data.Id;
          break;

        case "bleep":
          app.updateSwifts(data.blob);
          app.needsUpdate = true;
          break;

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
  svvarm,
  updateSwarm,
}
