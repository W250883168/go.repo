(function(root, factory) {
	if(typeof define === 'function' && define.amd) {
		define(['exports', 'echarts'], factory);
	} else if(typeof exports === 'object' && typeof exports.nodeName !== 'string') {
		factory(exports, require('echarts'));
	} else {
		factory({}, root.echarts);
	}
}(this, function(exports, echarts) {
	var log = function(msg) {
		if(typeof console !== 'undefined') {
			console && console.error && console.error(msg);
		}
	};
	if(!echarts) {
		log('ECharts is not Loaded');
		return;
	}
	if(!echarts.registerMap) {
		log('ECharts Map is not loaded');
		return;
	}
	
	echarts.registerMap('xinyun18', {
		"type": "FeatureCollection",
		"features": [{
			"type": "Feature",
			"id": "430102",
			"properties": {
				"name": "芙蓉区",
				//"cp": [112, 20],
				"childNum": 1
			},
			"geometry": {
				"type": "Polygon",
				"coordinates": [
					[
						[20, 20 * 0.75],
						[20, -20 * 0.75],
						[-20, -20 * 0.75],
						[-20, 20 * 0.75]
					]
				]
			}
		}]
	});
}));