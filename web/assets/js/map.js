// block: map — copy to web/assets/js/map.js. See map.md.
//
// Draws every [data-map] on the page with Leaflet, reading its points from
// the JSON the component wrote out. Markers always go on a canvas, since
// the home page has thousands; clusters and heat are opt-in per map and
// used only when their plugin loaded. Without Leaflet the fallback shows.

(function () {
	const TILES = 'https://tile.openstreetmap.org/{z}/{x}/{y}.png';
	const CREDIT = '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors';
	const still = window.matchMedia('(prefers-reduced-motion: reduce)').matches;

	function points(id) {
		const el = document.getElementById(id);
		try {
			return (el && JSON.parse(el.textContent)) || [];
		} catch (err) {
			return [];
		}
	}

	// Built from nodes, never an HTML string: every field here is GBIF's text.
	function popup(point) {
		const root = document.createElement('div');
		root.className = 'map__popup';
		root.appendChild(document.createElement('strong')).textContent = point.title;
		root.appendChild(document.createElement('span')).textContent =
			[point.date, point.place, point.count].filter(Boolean).join(' · ');
		if (point.stage) {
			root.appendChild(document.createElement('span')).textContent = point.stage;
		}
		if (point.anchor) {
			const link = root.appendChild(document.createElement('a'));
			link.href = '#' + point.anchor;
			link.textContent = 'Jump to record';
		}
		return root;
	}

	// Colours come from the element's CSS, so the palette lives in one place.
	function colour(el, prop, fallback) {
		return getComputedStyle(el).getPropertyValue(prop).trim() || fallback;
	}

	function init(el) {
		const data = points(el.dataset.map);
		if (!data.length || typeof L === 'undefined') {
			el.classList.add('is-unavailable');
			return;
		}

		const map = L.map(el, {
			preferCanvas: true,
			scrollWheelZoom: false,
			zoomAnimation: !still,
			fadeAnimation: !still,
			markerZoomAnimation: !still,
			worldCopyJump: true,
		});
		L.tileLayer(el.dataset.tiles || TILES, {
			attribution: el.dataset.attribution || CREDIT,
			maxZoom: 18,
			minZoom: 1,
		}).addTo(map);

		const latLngs = data.map((p) => [p.lat, p.lng]);
		let layer;
		if ('heat' in el.dataset && typeof L.heatLayer === 'function') {
			layer = L.heatLayer(latLngs, { radius: 18, blur: 15, minOpacity: 0.3 });
		} else {
			const fill = colour(el, '--map-marker', '#c2410c');
			const stroke = colour(el, '--map-marker-stroke', '#fff');
			const offFill = colour(el, '--map-marker-off', '#8a8a8a');
			layer = 'cluster' in el.dataset && typeof L.markerClusterGroup === 'function'
				? L.markerClusterGroup({ chunkedLoading: true, showCoverageOnHover: false })
				: L.featureGroup();
			data.forEach((p) => {
				L.circleMarker([p.lat, p.lng], {
					radius: p.off ? 5 : 7,
					weight: 2,
					color: stroke,
					fillColor: p.off ? offFill : fill,
					fillOpacity: p.off ? 0.6 : 0.9,
				})
					// A function, so the popup is only built when it opens.
					.bindPopup(() => popup(p))
					.addTo(layer);
			});
		}
		// Stages: a dashed line in running order, then a flag for each with its name.
		const stages = el.dataset.stages ? points(el.dataset.stages) : [];
		let frame = latLngs;
		if (stages.length) {
			const stageLatLngs = stages.map((s) => [s.lat, s.lng]);
			L.polyline(stageLatLngs, {
				color: colour(el, '--map-route', '#ffffff'),
				weight: 2,
				opacity: 0.8,
				dashArray: '6 6',
				interactive: false,
			}).addTo(map);
			stages.forEach((s) => {
				const kind = s.base ? 'base' : s.power ? 'power' : 'stage';
				const label = s.name + (s.power ? ' · Power Stage' : '') + (s.day && !s.base ? ' · ' + s.day : '');
				L.marker([s.lat, s.lng], {
					icon: L.divIcon({ className: 'stage-flag stage-flag--' + kind, iconSize: [16, 16] }),
					keyboard: false,
					title: label,
				})
					.bindTooltip(label, { direction: 'right', offset: [10, 0], className: 'stage-tip' })
					.addTo(map);
			});
			// Frame the stages and the sightings on the route, not every dot in Kenya.
			frame = stageLatLngs.concat(data.filter((p) => !p.off).map((p) => [p.lat, p.lng]));
		}
		layer.addTo(map);

		map.fitBounds(L.latLngBounds(frame), {
			padding: [24, 24],
			animate: false,
			maxZoom: Number(el.dataset.maxZoom) || 12,
		});

		// Wheel zoom waits for a click or focus, so the page scrolls past freely.
		map.on('focus', () => map.scrollWheelZoom.enable());
		map.on('blur', () => map.scrollWheelZoom.disable());
	}

	document.querySelectorAll('[data-map]').forEach(init);
})();
