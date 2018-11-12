import * as THREE from 'THREE';

import {sub, cross} from './math3.js';

// Return an instance of wavvey app
export class PlanetApp {
  constructor(params) {
    this.id = params.id;
    this.el = document.getElementById(this.id);
    this.app = {};
    this.width = this.el.offsetWidth;
    this.height = this.el.offsetHeight;
    this.state = params.state;
  }

  setup() {
    this.app = {
      view_angle: 67,
      aspect: this.width/this.height,
      near: 0.01,
      far: 200,
    };

    this.needsUpdate = false;

    // Scene
    this.scene = new THREE.Scene();

    // Renderer
    this.renderer = new THREE.WebGLRenderer({
      antialias: true,
      canvas: this.el,
    });

    this.renderer.setSize(this.width, this.height);
    this.renderer.setClearColor(0xFFFFFF, 1);
    this.renderer.setPixelRatio(1.5);

    document.body.appendChild(this.renderer.domElement);

    let light0 = new THREE.AmbientLight(0x777777);
    let light1 = new THREE.DirectionalLight( 0xcccccc, 0.4 );
    let light2 = new THREE.PointLight( 0xff0000, 0.5 );
    let light3 = new THREE.PointLight( 0x00ffff, 0.5 );
    light2.position.set(4.8, 10.0, -0.1);
    light3.position.set(4.0, 10.0, -0.0);

    this.scene.add(light0);
    // this.scene.add(light1);
    this.scene.add(light2);
    this.scene.add(light3);

    this.setupCamera();


    let a = 1.5;
    let [x, y, z] = [a-0.5, a, -a];
    this.resize(this.width, this.height);


    // Empty amount of ships
    this.ships = new Map();

  }

  // Return object containing all the necessary event handlers
  eventHandlers() {
    let self = this;
    let mouse = {x: 0.0, y: 0.0};

    return {
      resize: debounce(100, (ev) => {
        let size = Math.min(window.innerWidth, window.innerHeight);
        this.width = window.innerWidth;
        this.height = window.innerHeight;
        self.resize(this.width, this.height);
      }),
      move: debounce(10, (ev) => {
        mouse.x = ev.clientX;
        mouse.y = ev.clientY;
        // let u = 1*ev.clientX/window.innerWidth-1.0;
        // let v = 1*ev.clientX/window.innerWidth-1.0;
        // this.rgbPass.uniforms['amount'].value = Math.sqrt(u*u+v*v)/298.+0.002;
        let y = -1*(2*mouse.y/window.innerHeight - 1.0);
        let z = -1*(2*mouse.x/window.innerWidth - 1.0);
        //this.updatePosition(y, z);
      }),
    }
  }


  // Construct the torus
  setupWorld() {
  }

  coord(t, f, float) {
    let {radius, depth} = this.params;
    depth += float || 0.0;
    return [
      (radius+depth*Math.cos(f))*Math.cos(t),
      depth*Math.sin(f),
      (radius+depth*Math.cos(f))*Math.sin(t),
    ];
  }

  // Update world from a world object
  updateWorld(world) {

    this.params = {
      radius: world.radius,
      depth: world.thickness,
    };

    let geometry = new THREE.TorusGeometry(this.params.radius, this.params.depth, 16, 100);
    let material = new THREE.MeshPhongMaterial({
      color: 0xCCCCCC,
      emissive: 0x111111,
      specular: 0x444444,
      shininess: 90.0,
    });

    let mesh = new THREE.Mesh(geometry, material);
    mesh.rotation.x = Math.PI/2.0;
    mesh.rotation.z = Math.PI/2.0;
    this.scene.add(mesh);

    world.buildings.forEach((v) => {

      let height = 3*Math.random() + 0.1;
      height = 3;

      let {theta, fi} = v;
      // fi = 0.0;
      theta = Math.PI/2.0;


      let p = this.coord(theta, fi, 4.0);
      let q = this.coord(theta, fi, 6.0);
      let r = this.coord(theta+0.1, fi, 4.0);
      let b = sub(q, p);
      let c = sub(r, p);

      let f = new THREE.Vector3(b[0], b[1], b[2]);
      f.normalize();

      let o = new THREE.Mesh(
        new THREE.ConeGeometry(0.25, height),
        material,
      );

      console.log(r);
      o.position.set(p[0], p[1], p[2]);

      this.scene.add(o);

      /*
      this.scene.add(new THREE.ArrowHelper(
        f,
        o.position,
        4.3,
        0x000000,
      ));
      */

    });
  }

  /**
   * Setup Camera
   */
  setupCamera() {
    this.camera = new THREE.PerspectiveCamera(
      this.app.view_angle,
      this.app.aspect,
      this.app.near,
      this.app.far
    );

    // Camera
    this.camera.position.set(0.0, 30.0, 0.0);
    this.camera.lookAt(0, 0, 0);
  }

  getShipObject(id) {
    let o = this.scene.getObjectByName(id);

    if (o == undefined) {
      o = new THREE.Mesh(
        new THREE.BoxGeometry(0.5, 0.5, 0.5),
        new THREE.MeshBasicMaterial({color:0x000000}),
      );
      let a = Math.random()*2*Math.PI;
      let b = Math.random()*2*Math.PI;
      let c = Math.random()*2*Math.PI;
      o.rotation.set(a, b, c);
    }

    return o;
  }

  // Update positions of ships from a map
  updateShips(ships) {
    for (let id in ships) {
      let s = ships[id];
      let o = this.getShipObject("SHIP:"+id);

      if (!this.ships.has(id)) {
        o.name = "SHIP:"+id;
        this.scene.add(o);
      }

      let [x, y, z] = this.coord(s.coord.theta, s.coord.fi, 1.0);

      o.position.set(x, y, z);
      this.ships.set(id, s);
    }
  }

  // Update
  update(params) {
    if (this.needsUpdate) {
      this.needsUpdate = false;
    }

    let t = new Date()/1000.0;
    let r = 40.0;

    let x = r*Math.cos(t);
    let y = 0.0;
    let z = r*Math.sin(t);

    y = 40.0;

    this.camera.position.set(x, y, z);
    this.camera.lookAt(0.0, 0.0, 0.0);
  }

  // Resize canvas and set camera straight
  resize(width, height) {
    this.width = width;
    this.height = height;
    this.app.aspect = this.width/this.height;
    this.camera.aspect = this.app.aspect;
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(this.width, this.height);
  }

  // Draw
  draw() {
    this.renderer.render(this.scene, this.camera);
  }
}
