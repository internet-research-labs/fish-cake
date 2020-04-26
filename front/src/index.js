import {PlanetApp} from './PlanetApp.js';
import {SwarmApp} from './SwarmApp.js';

console.log("M A T T Y D I G I T A L");
console.log("aka Johnny Mnemonic");


/**
 * Game State
 */
let STATE = {
  ships: [],
};


let socket = undefined;


let Gap = (function () {
  let LAST = +new Date();
  return function () {
    let LATEST = + new Date();
    let v = LATEST - LAST;
    LAST = LATEST;
    return v;
  };
}())


function setupWebsocket(path) {

  let url = "ws://" + location.host + "/" + path;

  console.log("Connecting to", url);

  socket = new WebSocket(url);

  socket.addEventListener("open", function (ev) {
  });

  socket.addEventListener("close", function (ev) {
  });

  socket.addEventListener("message", function (ev) {
    console.log(ev.data);
  });
}


function svvarm(id) {


  // setupWebsocket("whatever");
  // return;


  // I never pay attention


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


  let LAST_TICK = 0;
  let PROCESSING = false;

  // Message
  socket.addEventListener('message', function (ev) {

    if (PROCESSING) {
      console.log("Returning early because we're processing");
      return;
    }

    let data = JSON.parse(ev.data);
    let type = data.type || data.Type;

    // Switch on this
    switch (type) {
      case "yupdate":

        PROCESSING = true;

        if (data.blob.tick-LAST_TICK != 1) {
          console.log("TICK MISMATCH:", data.blob.tick, LAST_TICK);
        }

        LAST_TICK = data.blob.tick;
        app.updateSwifts(data.blob.map);
        app.update();
        app.draw();

        PROCESSING = false;
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

function nvm(id) {

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

      if (PROCESSING) {
        console.log("Exiting early because we're still processing");
        return;
      }

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
    // requestAnimationFrame(loop);
  }());

  // Things that you might want later
  return [
    app,
  ];
}

function lezgo(id) {
}

export {
  lezgo,
  svvarm,
  updateSwarm,
}
