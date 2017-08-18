'use strict';

/**
 * 智慧教室设备管理系统
 * 包含：
 * 1、设备管理模块
 * 2、中控面板
 * 3、基础数据
 */


var config = {
	//HttpUrl: "http://192.168.0.201:8050",
	HttpUrl: "",
	loginuser: null,
	GetUser: function() {
		if(this.loginuser == null) {
			if(localStorage.getItem("LoginUser") != null && localStorage.getItem("LoginUser") != "undefined" && localStorage.getItem("LoginUser") != "") {
				this.loginuser = jQuery.parseJSON(localStorage.getItem("LoginUser"));
			} else {
				window.location.href = "/web2/html/login.html";
			}
		}
		return this.loginuser;
	},
	
	//   config  templates app.js
	
	zkmb_config:{
		//   教室
		classroomId: "1",
		//   班级
		classId: 0,
		//   用户ID
		Uid:'0',
		//   用于页面间传递某个设备
		item: {},
		//   类型
		Type: "classroom",
		//   当前班级课程章节中间ID
		Ccccid:"",
		//   图片路径
		deviceImg:"/web/upfile/device/",
		//   server
		coapServer:"http://192.168.0.201:8090",
		//   
		immediatelyRefreshTime:50,
		//
		fixRefreshTime:3000,
		//   中控点到论询延迟时间
		zkmbRefreshTime:5000,
		//   中控上下课轮询时间
		zkmbInOutRefreshTime:60000
	},
	//   设备监控刷新时间
	sbjkRefreshTime:5000,
    //   教室监控刷新时间
    jsjkRefreshTime:10000,
    //   实时出勤刷新时间
    sscqRefreshTime:5000,
    //   课表时间对象
    dateTable:{
		class_01: "08:00-08:45",
		class_02: "09:00-09:45",
		class_03: "10:00-10:45",
		class_04: "11:00-11:45",
		class_05: "12:00-12:45",
		class_06: "13:00-13:45",
		class_07: "14:00-14:45",
		class_08: "15:00-15:45",
		class_09: "16:00-16:45",
		class_10: "17:00-17:45",
		class_11: "18:00-18:45"
	}
};


//   UId
config.zkmb_config.Uid = config.GetUser().Usersid.toString();

angular.module('app', [
	'ngAnimate', 'ngCookies', 'ngResource', 'ngSanitize', 'ngTouch', 'ngStorage', 'ui.router', 'ui.bootstrap',
	'ui.load', 'ui.jq', 'ui.validate', 'oc.lazyLoad', 'pascalprecht.translate'
]);
var app =
	angular.module('app')
	.config(
		['$controllerProvider', '$compileProvider', '$filterProvider', '$provide',
			function($controllerProvider, $compileProvider, $filterProvider, $provide) {
				// lazy controller, directive and service
				app.controller = $controllerProvider.register;
				app.directive = $compileProvider.directive;
				app.filter = $filterProvider.register;
				app.factory = $provide.factory;
				app.service = $provide.service;
				app.constant = $provide.constant;
				app.value = $provide.value;
			}
		]);

angular.module('app')
	.constant('JQ_CONFIG', {
		easyPieChart: ['../vendor/jquery/charts/easypiechart/jquery.easy-pie-chart.js'],
		sparkline: ['../vendor/jquery/charts/sparkline/jquery.sparkline.min.js'],
		plot: ['../vendor/jquery/charts/flot/jquery.flot.min.js',
			'../vendor/jquery/charts/flot/jquery.flot.resize.js',
			'../vendor/jquery/charts/flot/jquery.flot.tooltip.min.js',
			'../vendor/jquery/charts/flot/jquery.flot.spline.js',
			'../vendor/jquery/charts/flot/jquery.flot.orderBars.js',
			'../vendor/jquery/charts/flot/jquery.flot.pie.min.js'
		],
		slimScroll: ['../vendor/jquery/slimscroll/jquery.slimscroll.min.js'],
		sortable: ['../vendor/jquery/sortable/jquery.sortable.js'],
		nestable: ['../vendor/jquery/nestable/jquery.nestable.js',
			'../vendor/jquery/nestable/nestable.css'
		],
		filestyle: ['../vendor/jquery/file/bootstrap-filestyle.min.js'],
		slider: ['../vendor/jquery/slider/bootstrap-slider.js',
			'../vendor/jquery/slider/slider.css'
		],
		chosen: ['../vendor/jquery/chosen/chosen.jquery.min.js',
			'../vendor/jquery/chosen/chosen.css'
		],
		TouchSpin: ['../vendor/jquery/spinner/jquery.bootstrap-touchspin.min.js',
			'../vendor/jquery/spinner/jquery.bootstrap-touchspin.css'
		],
		wysiwyg: ['../vendor/jquery/wysiwyg/bootstrap-wysiwyg.js',
			'../vendor/jquery/wysiwyg/jquery.hotkeys.js'
		],
		dataTable: ['../vendor/jquery/datatables/jquery.dataTables.min.js',
			'../vendor/jquery/datatables/dataTables.bootstrap.js',
			'../vendor/jquery/datatables/dataTables.bootstrap.css'
		],
		vectorMap: ['../vendor/jquery/jvectormap/jquery-jvectormap.min.js',
			'../vendor/jquery/jvectormap/jquery-jvectormap-world-mill-en.js',
			'../vendor/jquery/jvectormap/jquery-jvectormap-us-aea-en.js',
			'../vendor/jquery/jvectormap/jquery-jvectormap.css'
		],
		footable: ['../vendor/jquery/footable/footable.all.min.js',
			'../vendor/jquery/footable/footable.core.css'
		]
	})
	.config(['$ocLazyLoadProvider', function($ocLazyLoadProvider) {
		$ocLazyLoadProvider.config({
			debug: false,
			events: true,
			modules: [{
				name: 'ngGrid',
				files: [
					'../vendor/modules/ng-grid/ng-grid.min.js',
					'../vendor/modules/ng-grid/ng-grid.min.css',
					'../vendor/modules/ng-grid/theme.css'
				]
			}, {
				name: 'ui.select',
				files: [
					'../vendor/modules/angular-ui-select/select.min.js',
					'../vendor/modules/angular-ui-select/select.min.css'
				]
			}, {
				name: 'angularFileUpload',
				files: [
					'../vendor/modules/angular-file-upload/angular-file-upload.min.js'
				]
			}, {
				name: 'ui.calendar',
				files: ['../vendor/modules/angular-ui-calendar/calendar.js']
			}, {
				name: 'ngImgCrop',
				files: [
					'../vendor/modules/ngImgCrop/ng-img-crop.js',
					'../vendor/modules/ngImgCrop/ng-img-crop.css'
				]
			}, {
				name: 'angularBootstrapNavTree',
				files: [
					'../vendor/modules/angular-bootstrap-nav-tree/abn_tree_directive.js',
					'../vendor/modules/angular-bootstrap-nav-tree/abn_tree.css'
				]
			}, {
				name: 'toaster',
				files: [
					'../vendor/modules/angularjs-toaster/toaster.js',
					'../vendor/modules/angularjs-toaster/toaster.css'
				]
			}, {
				name: 'textAngular',
				files: [
					'../vendor/modules/textAngular/textAngular-sanitize.min.js',
					'../vendor/modules/textAngular/textAngular.min.js'
				]
			}, {
				name: 'vr.directives.slider',
				files: [
					'../vendor/modules/angular-slider/angular-slider.min.js',
					'../vendor/modules/angular-slider/angular-slider.css'
				]
			}, {
				name: 'com.2fdevs.videogular',
				files: [
					'../vendor/modules/videogular/videogular.min.js'
				]
			}, {
				name: 'com.2fdevs.videogular.plugins.controls',
				files: [
					'../vendor/modules/videogular/plugins/controls.min.js'
				]
			}, {
				name: 'com.2fdevs.videogular.plugins.buffering',
				files: [
					'../vendor/modules/videogular/plugins/buffering.min.js'
				]
			}, {
				name: 'com.2fdevs.videogular.plugins.overlayplay',
				files: [
					'../vendor/modules/videogular/plugins/overlay-play.min.js'
				]
			}, {
				name: 'com.2fdevs.videogular.plugins.poster',
				files: [
					'../vendor/modules/videogular/plugins/poster.min.js'
				]
			}, {
				name: 'com.2fdevs.videogular.plugins.imaads',
				files: [
					'../vendor/modules/videogular/plugins/ima-ads.min.js'
				]
			}, {
				name: 'jeDate',
				files: [
					'../vendor/jquery/jedate/skin/jedate.css',
					'../vendor/jquery/jedate/jedate.js'
				]
			}, {
				name: 'zkmb_blockui',
				files: [
					'../vendor/jquery/blockUI/jquery.blockUI.js'
				]
			}]
		});
	}]);

angular.module('app')
	.run(
		['$rootScope', '$state', '$stateParams',
			function($rootScope, $state, $stateParams) {
				$rootScope.$state = $state;
				$rootScope.$stateParams = $stateParams;
				
				//    返回
	            $rootScope.$on("$stateChangeSuccess",  function(event, toState, toParams, fromState, fromParams) {
	                // to be used for back button //won't work when page is reloaded.
	                $rootScope.previousState_name = fromState.name;
	                $rootScope.previousState_params = fromParams;
	            });
	            //back button function called from back button's ng-click="back()" 返回
	            $rootScope.back = function() {
	                $state.go($rootScope.previousState_name,$rootScope.previousState_params);
	            };
			}
		]
	)
	.config(
		['$stateProvider', '$urlRouterProvider',
			function($stateProvider, $urlRouterProvider) {
				$urlRouterProvider.otherwise('/app/index');
				$stateProvider
					.state('app', {
						abstract: true,
						url: '/app',
						templateUrl: 'app.html',
						resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									return $ocLazyLoad.load(
										[
											//     弹窗控制器
											//'../js/controllers/modal.js'
										]
									);
								}
							]
						}
					})
					.state('app.index', {
						url: '/index',
						templateUrl: 'app_index.html',
						resolve: {
							deps: ['$ocLazyLoad',
								function($ocLazyLoad) {
									//return $ocLazyLoad.load(['../js/controllers/chart.js']);
								}
							]
						}
					})

			}
		]
	);