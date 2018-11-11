import * as THREE from 'THREE';

// Return an instance of wavvey app
export class PlanetApp {
  constructor(params) {
    this.id = params.id;
    this.el = document.getElementById(this.id);
    this.app = {};
    this.width = this.el.offsetWidth;
    this.height = this.el.offsetHeight;
  }

  setup() {
    this.app = {
      view_angle: 50,
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
    this.setupWorld();


    let a = 1.5;
    let [x, y, z] = [a-0.5, a, -a];
    this.resize(this.width, this.height);
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


  setupWorld() {
    let geometry = new THREE.CylinderGeometry(1.0, 1.0, 16, 100);
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

    this.camera.position.set(0.0, 0.0, 40.0);
    this.camera.lookAt(0, 0, 0);
  }

  /**
   * Update
   */
  update(params) {
    if (this.needsUpdate) {
      this.needsUpdate = false;
    }
  }

  resize(width, height) {
    this.width = width;
    this.height = height;
    this.app.aspect = this.width/this.height;
    this.camera.aspect = this.app.aspect;
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(this.width, this.height);
    // this.composer.setSize(this.width, this.height);
  }

  draw() {
    // this.composer.render(1.05);
    this.renderer.render(this.scene, this.camera);
  }
}
