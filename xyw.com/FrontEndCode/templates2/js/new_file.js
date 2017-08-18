'use strict';

/* Controllers */

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
		])
	.config(['$translateProvider', function($translateProvider) {
		// Register a loader for the static files
		// So, the module will search missing translation tables under the specified urls.
		// Those urls are [prefix][langKey][suffix].
		$translateProvider.useStaticFilesLoader({
			prefix: '../l10n/',
			suffix: '.js'
		});
		// Tell the module what language to use by default
		$translateProvider.preferredLanguage('zn');
		// Tell the module to store the language in the local storage
		$translateProvider.useLocalStorage();
	}]);

angular.module('app')
	.controller('AppCtrl', ['$scope', '$translate', '$localStorage', '$window', 'httpService', '$rootScope', 'TipService', '$modal', '$state', '$interval','toaster',
		function($scope, $translate, $localStorage, $window, httpService, $rootScope, TipService, $modal, $state, $interval,toaster) {
			// add 'ie' classes to html
			var isIE = !!navigator.userAgent.match(/MSIE/i);
			isIE && angular.element($window.document.body).addClass('ie');
			isSmartDevice($window) && angular.element($window.document.body).addClass('smart');

			//  公共部份控制初始设置   true显示 ，false 隐藏
			$rootScope.showcom = true;

			//   config

			//console.log($localStorage.settings);

			// config
			$scope.app = {
				//   网站标题
				title:'网站标题',
				//   版本信息
				version: 'v1.1.0-beta',
				// for chart colors
				color: {
					primary: '#7266ba',
					info: '#23b7e5',
					success: '#27c24c',
					warning: '#fad733',
					danger: '#f05050',
					light: '#e8eff0',
					dark: '#3a3f51',
					black: '#1c2b36'
				},
				settings: {
					themeID: 13,
					navbarHeaderColor: 'bg-xyw-logo',
					navbarCollapseColor: 'bg-xyw-h',
					asideColor: 'bg-xyw-nav',
					headerFixed: true,
					asideFixed: false,
					asideFolded: false,
					asideDock: false,
					container: false
				},
				//   服务器时间
				serverTime: null,
				//   登录信息
				login:null,
				//   登录信息
				GetUser: function() {
					if(this.login == null) {
						if($window.localStorage.getItem("LoginUser") != null && $window.localStorage.getItem("LoginUser") != "undefined" && $window.localStorage.getItem("LoginUser") != "") {
							this.login = jQuery.parseJSON($window.localStorage.getItem("LoginUser"));
						} else {
							window.location.href = "/web2/html/login.html";
						}
					}
					return this.login;
				},
				//   http服务器地址
				httpUrl:"",
				//   coap服务器地址
				coapUrl:"",
				//   图片服务器地址
				imgUrl:"",
				//   轮询时间
				reTime:{
					//   设备监控刷新时间
					sbjkRefreshTime:5000,
				    //   教室监控刷新时间
				    jsjkRefreshTime:10000,
				    //   实时出勤刷新时间
				    sscqRefreshTime:5000
				},
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
				},
				//     中控面板config
				zkmb_config:{
					//   教室
					classroomId: "1",
					//
					immediatelyRefreshTime:50,
					//
					fixRefreshTime:3000,
					//   中控点到论询延迟时间
					zkmbRefreshTime:5000,
					//   中控上下课轮询时间
					zkmbInOutRefreshTime:60000
				}
			};

			//   是否登录
			if(angular.isDefined($scope.app.GetUser())) {
				$scope.app.login = $scope.app.GetUser();
			} else {
				window.location.href = "/web2/html/login.html";
			}

			//   登录信息
			$scope.app.login = $scope.app.GetUser();

			//    nav
			$scope.navlist = [];
			//    取模块列表
			var getapp = function(){
				var url = $scope.app.httpUrl + "/getapp";
				var data = {
					"Usersid": $scope.app.login.Usersid,
					"Rolestype": $scope.app.login.Rolestype,
					"Os": "PcWEB",
					"Token": $scope.app.login.Token
				};
				var promise = httpService.ajaxPost(url, data);
				promise.then(function(data) {
					if(data.Rcode == "1000") {
						$scope.navlist = data.Result;
					} else {
            toaster.pop('warning', data.Reason);
					}
				}, function(reason) {}, function(update) {});
			}


			// save settings to local storage
			if(angular.isDefined($localStorage.settings)) {
				$scope.app.settings = $localStorage.settings;
			} else {
				$localStorage.settings = $scope.app.settings;
			}
			$scope.$watch('app.settings', function() {
				if($scope.app.settings.asideDock && $scope.app.settings.asideFixed) {
					// aside dock and fixed must set the header fixed.
					$scope.app.settings.headerFixed = true;
				}
				// save to local storage
				$localStorage.settings = $scope.app.settings;
			}, true);

			function isSmartDevice($window) {
				// Adapted from http://www.detectmobilebrowsers.com
				var ua = $window['navigator']['userAgent'] || $window['navigator']['vendor'] || $window['opera'];
				// Checks for iOs, Android, Blackberry, Opera Mini, and Windows mobile devices
				return(/iPhone|iPod|iPad|Silk|Android|BlackBerry|Opera Mini|IEMobile/).test(ua);
			}

			//    退出
			$scope.loginout = function() {
				localStorage.removeItem("LoginUser");
				window.location.href = "/web2/html/login.html";
			};

			//  是否有子目录    return : true 为干节点，false 为叶子节点
			$scope.haveSon = function(id) {
				var haves = false;
				for(var i = 0; i < $scope.navlist.length; i++) {
					if($scope.navlist[i].Superiormoduleid == id) {
						haves = true;
						return haves;
					}
				}
				return haves;
			}

			//  提示  弹窗
			$scope.tipService = TipService;

			//   左侧导航跳转
			$scope.navClick = function(item) {
				if(!$scope.haveSon(item.Id)) {
					$state.go("app." + item.Modulecode, {}, {
						reload: true
					});
				}
			};

			/**
			 * 取服务器时间
			 */
			$scope.getServerTimeTimer = null;
			$scope.getServerTime = function() {
				var url = $scope.app.httpUrl + "/getServerTime";
				var data = {
					"Usersid": $scope.app.login.Usersid,
					"Rolestype": $scope.app.login.Rolestype,
					"Os": "WEB",
					"Token": $scope.app.login.Token
				};
				var promise = httpService.ajaxPost(url, data);
				promise.then(function(data) {
					console.log('取服务器时间', data);
					if(data.Rcode == "1000") {
						//  服务器时间
						$interval.cancel($scope.getServerTimeTimer);
						$scope.getServerTimeTimer = null;
						$scope.app.serverTime = data.Result;
						$scope.getServerTimeTimer = $interval(function() {
							$scope.app.serverTime += 1;
						}, 1000);
					} else {
						//   本地时间
						$interval.cancel($scope.getServerTimeTimer);
						$scope.getServerTimeTimer = null;
						$scope.app.serverTime = Date.parse(new Date()) / 1000;
						$scope.getServerTimeTimer = $interval(function() {
							$scope.app.serverTime += 1;
						}, 1000);

						//$scope.tipService.setMessage(data.Reason);
					}
				}, function(reason) {}, function(update) {});
			}

			//   轮询
			$scope.setglobaldata = {
				timedatalist: [], //定时器数组
				timer: { //定时器模板对象
					interval: null,
					Key: "", //定义的名称
					keyctrl: "", //定义所属的控制器
					fnStopAutoRefresh: function() {}, //定义开关的关闭
					fnAutoRefresh: function() {}, //定义开关的打开
					fnStopAutoRefreshfn: function(tm, fn) {}, //定义开关的关闭方法
					fnAutoRefreshfn: function(tm) {
						//console.log(tm.keyctrl);
						console.log('$state.current.name',$state.current.name);
						if(tm.keyctrl != $state.current.name) {
							tm.fnStopAutoRefresh();
						} else {
							if(tm.interval == null) {
								tm.fnAutoRefresh();
							}
						}
					}
				},
				addtimer: function(t) { //将数据加入到定时器数组
					var isadd = true;
					//console.log(t);
					for(var i = 0; i < this.timedatalist.length; i++) {
						if(this.timedatalist[i].Key == t.key) {
							this.timedatalist[i].fnStopAutoRefresh(); //先关闭定时器
							this.timedatalist.splice(i, 1); //移除对象
						}
					}
					if(isadd) {
						this.timer = t;
						this.timedatalist.push(this.timer);
					}
				},
				gettimer: function(key) { //获取对象
					for(var i = 0; i < this.timedatalist.length; i++) {
						if(this.timedatalist[i].Key == key) {
							this.timer = this.timedatalist[i];
							break;
						}
					}
					return angular.copy(this.timer);
				}
			};
			//console.log($state.current.name);
			//   监听离开页面取消定时器
			$rootScope.$on('$stateChangeSuccess',
				function(event, toState, toParams, fromState, fromParams) {
					//console.log("监听离开页面取消定时器")
					//console.log(toState.name);
					//console.log(fromState.name);
					for(var indextm = 0; indextm < $scope.setglobaldata.timedatalist.length; indextm++) {
						if($scope.setglobaldata.timedatalist[indextm].keyctrl == toState.name) {
							$scope.setglobaldata.timedatalist[indextm].fnAutoRefresh();
						} else {
							$scope.setglobaldata.timedatalist[indextm].fnStopAutoRefresh();
						}
					}
				}
			);

			var run = function() {
				getapp();
				$scope.getServerTime();
			}
			run();

		}
	]);
