// A gravel rally wheel spinning in the banner: a treaded tyre, a white six-spoke rim,
// a red hub and a green brake caliper. Low poly, so it stays light on phones.

import { mount } from '/assets/js/threejs.js';

mount('rally-wheel', function ({ THREE, scene, camera, colour, pointer }) {
	const tilt = new THREE.Group();
	const wheel = new THREE.Group();
	tilt.add(wheel);
	scene.add(tilt);

	const rubber = new THREE.MeshStandardMaterial({ color: 0x18181a, roughness: 0.95 });
	const rimPaint = new THREE.MeshStandardMaterial({ color: colour('--three-rim', '#ffffff'), roughness: 0.35, metalness: 0.2 });
	const recess = new THREE.MeshStandardMaterial({ color: 0x2a2a2e, roughness: 0.6, metalness: 0.4 });
	const red = new THREE.MeshStandardMaterial({ color: colour('--three-hub', '#bb0000'), roughness: 0.4 });

	// The axle runs along z, so the wheel faces the camera and turns about rotation.z.
	const tyre = new THREE.Mesh(new THREE.TorusGeometry(1.28, 0.4, 14, 48), rubber);
	wheel.add(tyre);

	// Gravel tread: chunky blocks round the outside, two rows staggered.
	const blocks = 40;
	const tread = new THREE.InstancedMesh(new THREE.BoxGeometry(0.16, 0.12, 0.26), rubber, blocks);
	const place = new THREE.Object3D();
	for (let i = 0; i < blocks; i++) {
		const a = (i / blocks) * Math.PI * 2;
		place.position.set(Math.cos(a) * 1.68, Math.sin(a) * 1.68, i % 2 ? 0.15 : -0.15);
		place.rotation.set(0, 0, a);
		place.updateMatrix();
		tread.setMatrixAt(i, place.matrix);
	}
	wheel.add(tread);

	// Cylinders stand on y by default; turn them to lie along the axle.
	function disc(radius, depth, material, z, segments) {
		const mesh = new THREE.Mesh(new THREE.CylinderGeometry(radius, radius, depth, segments || 32), material);
		mesh.rotation.x = Math.PI / 2;
		mesh.position.z = z;
		return mesh;
	}

	wheel.add(disc(0.95, 0.28, recess, -0.02));
	// The rim's lip, a thin white ring at the front.
	const lip = new THREE.Mesh(new THREE.TorusGeometry(0.93, 0.06, 8, 40), rimPaint);
	lip.position.z = 0.16;
	wheel.add(lip);

	for (let i = 0; i < 6; i++) {
		const a = (i / 6) * Math.PI * 2;
		const spoke = new THREE.Mesh(new THREE.BoxGeometry(0.2, 0.75, 0.1), rimPaint);
		spoke.position.set(Math.cos(a) * 0.52, Math.sin(a) * 0.52, 0.17);
		spoke.rotation.z = a - Math.PI / 2;
		wheel.add(spoke);
	}

	wheel.add(disc(0.24, 0.18, red, 0.2, 6));
	wheel.add(disc(0.1, 0.12, rubber, 0.3, 6));

	// The caliper sits still behind the rim while the wheel turns.
	const caliper = new THREE.Mesh(
		new THREE.BoxGeometry(0.34, 0.6, 0.22),
		new THREE.MeshStandardMaterial({ color: colour('--three-caliper', '#006600'), roughness: 0.5 }),
	);
	caliper.position.set(0.62, 0.32, -0.1);
	caliper.rotation.z = -0.5;
	tilt.add(caliper);

	scene.add(new THREE.AmbientLight(0xffffff, 0.7));
	const key = new THREE.DirectionalLight(0xffffff, 2.4);
	key.position.set(3, 4, 5);
	scene.add(key);
	const dust = new THREE.DirectionalLight(colour('--three-dust', '#e89a5f'), 2.5);
	dust.position.set(-4, -1, -2);
	scene.add(dust);

	// A three-quarter view, so it reads as a wheel even when it's standing still.
	tilt.rotation.set(0.12, -0.55, 0);
	camera.position.z = 5.4;

	return function (time) {
		wheel.rotation.z = -time * 2.4;
		tilt.rotation.y = -0.55 + pointer.x * 0.25;
		tilt.rotation.x = 0.12 - pointer.y * 0.15;
		tilt.position.y = Math.sin(time * 9) * 0.015;
	};
});
