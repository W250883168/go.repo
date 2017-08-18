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
			// config
			$scope.app = {
				name: config.GetUser().Truename,
        title: '智慧教室',
				headimg: config.loginuser.Userheadimg,
				version: 'v1.1.0-beta',
				// for chart colors
				color: {
		          primary: '#7266ba',
		          info:    '#2ea2f8',
		          success: '#27c24c',
		          warning: '#fad733',
		          danger:  '#f05050',
		          light:   '#e8eff0',
		          dark:    '#3a3f51',
		          black:   '#1b2431'
		        },
		        settings: {
		          themeID: 1,
		          navbarHeaderColor: 'bg-black b-r-b b-b-b',
		          navbarCollapseColor: 'bg-black',
		          asideColor: 'bg-black b-r-b',
              	  bodyColor: 'bg-black-body app-black',
		          footerColor: 'bg-black b-t-b',
		          headerFixed: true,
		          asideFixed: false,
		          asideFolded: false,
		          asideDock: false,
		          container: false
		        },
				//   服务器时间
				serverTime: null,
				//   登录信息
				'login': {}
			};
			
			//   设备图片路径
			$scope.deviceImg = config.zkmb_config.deviceImg;

			//   登录信息
			$scope.app.login = config.GetUser();

			$scope.navlist = [];
			var url = config.HttpUrl + "/getapp";

			var data = {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "PcWEB",
				"Token": config.GetUser().Token
			};
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				if(data.Rcode == "1000") {
					$scope.navlist = data.Result;
					console.log($scope.navlist);
				} else {
          			toaster.pop('warning',data.Reason);
				}
			}, function(reason) {}, function(update) {});
			if(angular.isDefined(config.GetUser())) {
				config.loginuser = config.GetUser();
			} else {
				window.location.href = "/web2/html/login.html";
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
				var url = config.HttpUrl + "/getServerTime";
				var data = {
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Os": "WEB",
					"Token": config.GetUser().Token
				};
				var promise = httpService.ajaxPost(url, data);
				promise.then(function(data) {
					console.log('取服务器时间', data);
					if(data.Rcode == "1000") {
						$interval.cancel($scope.getServerTimeTimer);
						$scope.getServerTimeTimer = null;
						$scope.app.serverTime = data.Result;
						$scope.getServerTimeTimer = $interval(function() {
							$scope.app.serverTime += 1;
						}, 1000);
					} else {
						//$scope.tipService.setMessage(data.Reason);
					}
				}, function(reason) {}, function(update) {});
			}

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

			//  =============================   定义枚举数据  ===================================
			/**
			 * 设备模板
			 */
			$scope.enumDeviceModel = {
				//  灯
				'lamp':{
					'title':'灯',
					'PageFileName':'lamp.html',
					'ImgFileName':'lamp2.png',
					'ImgFileName2':'lamp2.png'
				},
				//  空调
				'kt':{
					'title':'空调',
					'PageFileName':'kt.html',
					'ImgFileName':'kt2.png',
					'ImgFileName2':'kt2.png'
				},
				//  投影仪
				'projector':{
					'title':'投影仪',
					'PageFileName':'projector.html',
					'ImgFileName':'projector2.png',
					'ImgFileName2':'projector2.png'
				},
				//   VGA切换器
				'vgatoggle':{
					'title':'VGA切换器',
					'PageFileName':'vgatoggle.html',
					'ImgFileName':'vgatoggle2.png',
					'ImgFileName2':'vgatoggle2.png'
				},
				//   pjlink
				'pjlink':{
					'title':'pjlink',
					'PageFileName':'pjlink.html',
					'ImgFileName':'pjlink2.png',
					'ImgFileName2':'pjlink2.png'
				},
				//   版本切换
				'upgrade':{
					'title':'版本切换',
					'PageFileName':'upgrade.html',
					'ImgFileName':'vgatoggle2.png',
					'ImgFileName2':'vgatoggle2.png'
				}
			}


			//  =============================   /定义枚举数据  ===================================

			//   ================================  alert 提示  =====================================
//			error: 'toast-error',
//	        info: 'toast-info',
//	        wait: 'toast-wait',
//	        success: 'toast-success',
//	        warning: 'toast-warning'

			$scope.toaster = {
		        type: 'success',
		        title: 'Title',
		        text: 'Message'
		    };
		    $scope.pop = function(){
		        toaster.pop($scope.toaster.type, $scope.toaster.title, $scope.toaster.text);
		    };
			//   ================================  /alert 提示  =====================================




			var run = function() {
				$scope.getServerTime();
			}
			run();

		}
	]);

app.controller('appIndexContr', ['$scope', 'httpService','$filter','toaster', function($scope, httpService,$filter,toaster) {
	console.log("首页")
	$scope.myChart1={};
	$scope.myChart1.Chart={};
	$scope.myChart1.xAxis={};
	$scope.myChart1.xAxis.data=['整体', '大一', '大二', '大三', '大四'];
	$scope.myChart1.series={};
	$scope.myChart1.series.data=[];
	$scope.myChart1.dateitem=null;
	$scope.myChart1.gradeitem={};
	$scope.myChart1.CollegeItem={};
	$scope.myChart1.MajorItem={};

	$scope.myChart2={};
	$scope.myChart2.Chart={};
	$scope.myChart2.xAxis={};
	$scope.myChart2.xAxis.data=[];
	$scope.myChart2.xAxis.config=[];
	$scope.myChart2.series={};
	$scope.myChart2.series.data=[];
	$scope.myChart2.dateitem=30;
	$scope.myChart2.gradeitem={};
	$scope.myChart2.CollegeItem={};
	$scope.myChart2.MajorItem={};

	$scope.myChart3={};
	$scope.myChart3.dateitem=30;
	$scope.myChart3.gradeitem={};
	$scope.myChart3.CollegeItem={};
	$scope.myChart3.MajorItem={};
	$scope.myChart3.Chart={};
	$scope.myChart3.xAxis={};
	$scope.myChart3.xAxis.data=[];
	$scope.myChart3.xAxis.config=[];
	$scope.myChart3.series={};
	$scope.myChart3.series.data=[];

	$scope.myChart4={};
	$scope.myChart4.dateitem=30;
	$scope.myChart4.gradeitem={};
	$scope.myChart4.CollegeItem={};
	$scope.myChart4.MajorItem={};
	$scope.myChart4.Chart={};
	$scope.myChart4.xAxis={};
	$scope.myChart4.xAxis.data=[];
	$scope.myChart4.xAxis.config=[];
	$scope.myChart4.series={};
	$scope.myChart4.series.data=[];

	$scope.myChart5={};
	$scope.myChart5.dateitem=30;
	$scope.myChart5.gradeitem={};
	$scope.myChart5.CollegeItem={};
	$scope.myChart5.MajorItem={};
	$scope.myChart5.Chart={};
	$scope.myChart5.xAxis={};
	$scope.myChart5.xAxis.data=[];
	$scope.myChart5.xAxis.config=[];
	$scope.myChart5.series={};
	$scope.myChart5.series.data=[];

	$scope.College={};
	$scope.College.selectItem={};
	$scope.College.Collegelist=[];
	$scope.Major={};
	$scope.Major.selectItem={};
	$scope.Major.Majorlist=[];
	//	时间
	$scope.NowDate = '';
	//	早晨
	$scope.Hours = '';
	
	var t = null;
	//开始执行
 	t = setTimeout(time,1000);
	//	获取当前时间
	var time = function () {
		//清除定时器
	 	clearTimeout(t);
		var now = new Date();
		//	判断时间
		var hou = now.getHours();
		if (hou >= 6 && hou < 8) {
			$scope.Hours = '早上好';
		} else if (hou >= 8 && hou < 12) {
			$scope.Hours = '上午好';
		} else if (hou >=12 && hou < 19) {
			$scope.Hours = '下午好';
		} else {
			$scope.Hours = '晚上好';
		}
		$scope.NowDate = $filter('date')(now,'yyyy年MM月dd日HH:mm:ss');
		//设定定时器，循环执行
		t = setTimeout(time,1000);
	}
	
	var Init_load=function()//初始化加载相关设置数据
	{
		var url = config.HttpUrl+"/basicset/getall";
        var promise =httpService.ajaxGet(url,null);
        promise.then(function (data) {
        	console.log(data.Result[4])
            if(data.Rcode=="1000"){
                $scope.College.Collegelist=data.Result[3];
                $scope.Major.Majorlist=data.Result[4];
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
        $scope.Load_ValChart1();
        $scope.Load_ValChart2();
        $scope.Load_ValChart3();
        $scope.Load_ValChart4();
        $scope.Load_ValChart5();
	};

	var HandleChar1=function(d){//整体出勤分析
    $scope.myChart1.series.data=[];
    var count=0;
    for(var i=0;i<d.length;i++){
      count=count+d[i].Analysisvalue;
      $scope.myChart1.series.data.push(Math.round(d[i].Analysisvalue.toFixed(2)*10000)/100);
    }
    $scope.myChart1.series.data.push(Math.round( (count/d.length).toFixed(2) * 10000 ) /100);
    $scope.myChart1.series.data=$scope.myChart1.series.data.reverse();
    $scope.htmlReady1();
	};
	var HandleChar2=function(d){
		$scope.myChart2.series.data=[];
		$scope.myChart2.xAxis.data=[];
		var count=0;
		for(var i=0;i<d.length;i++){
			count=count+d[i].Analysisvalue;
			$scope.myChart2.series.data.push(Math.round(d[i].Analysisvalue.toFixed(2)*10000) /100);
			$scope.myChart2.xAxis.data.push(d[i].Analysisname);
		}
		$scope.htmlReady2();
	};
	var HandleChar3=function(d){
		$scope.myChart3.series.data=[];
		$scope.myChart3.xAxis.data=[];
		var count=0;
		for(var i=0;i<d.length;i++){
			count=count+d[i].Analysisvalue;
			$scope.myChart3.series.data.push(Math.round(d[i].Analysisvalue.toFixed(2)*10000) / 100);
			$scope.myChart3.xAxis.data.push(d[i].Analysisname);
		}
		$scope.htmlReady3();
	};
	var HandleChar4=function(d){
		$scope.myChart4.series.data=[];
		$scope.myChart4.xAxis.data=[];
		var count=0;
		for(var i=0;i<d.length;i++){
			count=count+d[i].Analysisvalue;
			$scope.myChart4.series.data.push(Math.round(d[i].Analysisvalue.toFixed(2)*10000) / 100);
			$scope.myChart4.xAxis.data.push(d[i].Analysisname);
		}
		$scope.htmlReady4();
	};
	var HandleChar5=function(d){
		$scope.myChart5.series.data=[];
		$scope.myChart5.xAxis.data=[];
		var count=0;
		for(var i=0;i<d.length;i++){
			count=count+d[i].Analysisvalue;
			$scope.myChart5.series.data.push(Math.round(d[i].Analysisvalue.toFixed(2)*10000) / 100);
			$scope.myChart5.xAxis.data.push(d[i].Analysisname);
		}
		$scope.htmlReady5();
	};

	$scope.Load_ValChart1=function(item){//整体出勤分析
    if (!item) {
      return;
    } else {
      $scope.myChart1.dateitem = item.value;
    }
		var url = config.HttpUrl+"/curriculum/attendanceanalysis";
        var data={
        "Usersid": config.GetUser().Usersid,
        "Rolestype": config.GetUser().Rolestype,
        "Token": config.GetUser().Token,
        "Os": "WEB",
        "Dateint": Number($scope.myChart1.dateitem),
        "Gradeint": Number($scope.myChart1.gradeitem.value),
        "Majorid": Number($scope.myChart1.MajorItem.Id),
        "Collegeid":Number($scope.myChart1.CollegeItem.Id),
        "Curriculumsid":0,
        "Analysistype":0
      };
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
              if(data.Result!=null){
                HandleChar1(data.Result);
              }
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
	};
	$scope.Load_ValChart2=function(name,item){//学院出勤分析
    if (!item) {
      return;
    } else if (name == 'date') {
      $scope.myChart2.dateitem = item.value;
    } else if (name == 'grade') {
      $scope.myChart2.gradeitem = item.value;
    }
		var url = config.HttpUrl+"/curriculum/attendanceanalysis";
        var data={
          "Usersid": config.GetUser().Usersid,
          "Rolestype": config.GetUser().Rolestype,
          "Token": config.GetUser().Token,
          "Os": "WEB",
          "Dateint": Number($scope.myChart2.dateitem),
          "Gradeint": Number($scope.myChart2.gradeitem),
          "Majorid": Number($scope.myChart2.MajorItem.Id),
          "Collegeid":Number($scope.myChart2.CollegeItem.Id),
          "Curriculumsid":0,
          "Analysistype":1
        };
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
            	if(data.Result!=null){
            		HandleChar2(data.Result);
            	}
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
	};
	$scope.Load_ValChart3=function(name,item){//专业出勤分析
    if (!item) {
      return;
    } else if (name == 'date') {
      $scope.myChart3.dateitem = item.value;
    } else if (name == 'grade') {
      $scope.myChart3.gradeitem = item.value;
    }
		var url = config.HttpUrl+"/curriculum/attendanceanalysis";
        var data={
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Dateint": Number($scope.myChart3.dateitem),
			"Gradeint": Number($scope.myChart3.gradeitem),
			"Majorid": Number($scope.myChart3.MajorItem.Id),
			"Collegeid":Number($scope.myChart3.CollegeItem),
			"Curriculumsid":0,
			"Analysistype":2
		};
    var promise =httpService.ajaxPost(url,data);
    promise.then(function (data) {
        if(data.Rcode=="1000"){
          if(data.Result!=null){
          HandleChar3(data.Result);
          }
        }else{
          toaster.pop('warning',data.Reason);
        }
    }, function (reason) {}, function (update) {});
	};
	$scope.Load_ValChart4=function(name,item){//班级出勤分析
    if (!item) {
      return;
    } else if (name == 'date') {
      $scope.myChart4.dateitem = item.value;
    } else if (name == 'grade') {
      $scope.myChart4.gradeitem = item.value;
    } else if (name == 0) {
      $scope.myChart4.MajorItem = item.Id;
    } else if (name == 1) {
      $scope.myChart4.CollegeItem = item.Id;
    }
		var url = config.HttpUrl+"/curriculum/attendanceanalysis";
        var data={
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Dateint": Number($scope.myChart4.dateitem),
			"Gradeint": Number($scope.myChart4.gradeitem),
			"Majorid": Number($scope.myChart4.MajorItem),
			"Collegeid":Number($scope.myChart4.CollegeItem),
			"Curriculumsid":0,
			"Analysistype":3
		};
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
            	if(data.Result!=null){
            		HandleChar4(data.Result);
            	}
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
	};
	$scope.Load_ValChart5=function(name,item){//课程出勤分析
    if (!item) {
      return;
    } else if (name == 'date') {
      $scope.myChart5.dateitem = item.value;
    } else if (name == 'grade') {
      $scope.myChart5.gradeitem = item.value;
    } else if (name == 0) {
      $scope.myChart5.MajorItem = item.Id;
    } else if (name == 1) {
      $scope.myChart5.CollegeItem = item.Id;
    }
		var url = config.HttpUrl+"/curriculum/attendanceanalysis";
        var data={
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Dateint": Number($scope.myChart5.dateitem),
			"Gradeint": Number($scope.myChart5.gradeitem),
			"Majorid": Number($scope.myChart5.MajorItem),
			"Collegeid":Number($scope.myChart5.CollegeItem),
			"Curriculumsid":0,
			"Analysistype":4
		};
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
            	if(data.Result!=null){
            		HandleChar5(data.Result);
            	}
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
	};

	$scope.htmlReady1 = function() {
//		// 基于准备好的dom，初始化echarts实例
		$scope.myChart1.Chart = echarts.init(document.getElementById('sbfx_kt1'));
		// 指定图表的配置项和数据
		var option = {
      fn: true,
			tooltip: {
				trigger: 'item'
			},
			toolbox: {
				show: false,
				feature: {
					dataView: {
						show: true,
						readOnly: false
					},
					restore: {
						show: true
					},
					saveAsImage: {
						show: true
					}
				}
			},
			calculable: true,
      grid: {
        borderWidth: 0,
        bottom:60,
        top:'20%',
        y: 80,
        y2: 60
      },
			xAxis: [{
				type: 'category',
				show: true,
				data: $scope.myChart1.xAxis.data,
				axisLabel: {margin:10,textStyle: {color: '#7f8fa4',fontSize: 14,fontWeight: ''}},
				axisLine: {show: false},
				axisTick:{show:false}
			}],
			yAxis: [{type: 'value',show: false}],
			series: [{
				//name: '空调',
				type: 'bar',
        barWidth: 18,
				label:{normal:{show:true,position:'top',formatter:'{c}' + '%',textStyle:{color:'#7f8fa4',fontSize:"14"}}
				},
				itemStyle: {
					normal: {
						color: function(params) {
							var colorList = ['#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0'];
							return colorList[params.dataIndex]
						},
            barBorderRadius:[20,20,0,0]
					}
				},
				data: [99, 21, 10, 4, 12]//$scope.myChart1.xAxis.data
			}]
		};
		// 使用刚指定的配置项和数据显示图表。
		$scope.myChart1.Chart.setOption(option);
	}
	$scope.htmlReady2 = function() {
		// 基于准备好的dom，初始化echarts实例
		$scope.myChart2.Chart = echarts.init(document.getElementById('sbfx_kt2'));
//		// 指定图表的配置项和数据
		var option2 = {
      fn: true,
			tooltip: {
				trigger: 'item'
			},
			toolbox: {
				show: false,
				feature: {dataView: {show: true,readOnly: false},
					restore: {show: true},
					saveAsImage: {show: true}
				}
			},
			calculable: true,
      grid: {
        borderWidth: 0,
        bottom:60,
        top:'20%',
        y: 80,
        y2: 60
      },
			xAxis: [{
				type: 'category',
				show: true,
				data: $scope.myChart2.xAxis.data,//['本部校区', '新校区', '北校区', '东校区', '本校区','南校区', '老校区', '中部校区', '东南校区', '本南校区'],
				axisLabel: {
					margin:20,
					textStyle: {
						color: '#7f8fa4',
						fontSize: 14,
						fontWeight: ''
					}

				},
				axisLine: {
					show: false
				},
				axisTick:{
				    show:false
				}
			}],
			yAxis: [{
				type: 'value',
				show: false
			}],
			series: [{
				name: '',
				type: 'bar',
				barWidth: 18,
				label:{
				    normal:{
				        show:true,
				        position:'top',
				        formatter:'{c}' + '%',
				        textStyle:{
				            color:'#7f8fa4',
				            fontSize:"14"
				        }
				    }
				},
				itemStyle: {
					normal: {
						color: function(params) {
							// build a color map as your need.
							var colorList = ['#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0','#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0'];
							return colorList[params.dataIndex]
						},
            barBorderRadius:[50,50,0,0]
          }
				},
				data: $scope.myChart2.xAxis.data //[124,156,22,6,9,45]
			}]
		};
		$scope.myChart2.Chart.setOption(option2);
	}
  $scope.htmlReady3 = function() {
    // 基于准备好的dom，初始化echarts实例
    $scope.myChart3.Chart = echarts.init(document.getElementById('sbfx_kt3'));
//		// 指定图表的配置项和数据
    var option3 = {
      fn: true,
      tooltip: {
        trigger: 'item'
      },
      toolbox: {
        show: false,
        feature: {dataView: {show: true,readOnly: false},
          restore: {show: true},
          saveAsImage: {show: true}
        }
      },
      calculable: true,
      grid: {
        borderWidth: 0,
        bottom:60,
        top:'20%',
        y: 80,
        y2: 60
      },
      xAxis: [{
        type: 'category',
        show: true,
        data: $scope.myChart3.xAxis.data,//['英语系', '法语系', '日语系', '西班牙语系'],
        axisLabel: {
          margin:20,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          }

        },
        axisLine: {
          show: false
        },
        axisTick:{
          show:false
        }
      }],
      yAxis: [{
        type: 'value',
        show: false
      }],
      series: [{
        name: '',
        type: 'bar',
        barWidth: 18,
        label:{
          normal:{
            show:true,
            position:'top',
            formatter:'{c}' + '%',
            textStyle:{
              color:'#7f8fa4',
              fontSize:"14"
            }
          }
        },
        itemStyle: {
          normal: {
            color: function(params) {
              // build a color map as your need.
              var colorList = ['#2297F0', '#2297F0', '#2297F0', '#2297F0'];
              return colorList[params.dataIndex]
            },
            barBorderRadius:[50,50,0,0]
          }
        },
        data:$scope.myChart3.series.data//[99, 21, 10, 4]
      }]
    };
    $scope.myChart3.Chart.setOption(option3);
  }
  $scope.htmlReady4 = function() {
    // 基于准备好的dom，初始化echarts实例
    $scope.myChart4.Chart = echarts.init(document.getElementById('sbfx_kt4'));
//		// 指定图表的配置项和数据
    var option4 = {
      fn: true,
      tooltip: {
        trigger: 'item'
      },
      toolbox: {
        show: false,
        feature: {dataView: {show: true,readOnly: false},
          restore: {show: true},
          saveAsImage: {show: true}
        }
      },
      calculable: true,
      grid: {
        borderWidth: 0,
        bottom:60,
        top:'20%',
        y: 80,
        y2: 60
      },
      xAxis: [{
        type: 'category',
        show: true,
        data: $scope.myChart4.xAxis.data,//['英语系1班', '英语系2班', '英语系3班', '英语系4班', '英语系5班','英语系6班', '英语系7班', '英语系8班', '英语系9班', '英语系10班'],
        axisLabel: {
          margin:20,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          }

        },
        axisLine: {
          show: false
        },
        axisTick:{
          show:false
        }
      }],
      yAxis: [{
        type: 'value',
        show: false
      }],
      series: [{
        name: '',
        type: 'bar',
        barWidth: 18,
        label:{
          normal:{
            show:true,
            position:'top',
            formatter:'{c}' + '%',
            textStyle:{
              color:'#7f8fa4',
              fontSize:"14"
            }
          }
        },
        itemStyle: {
          normal: {
            color: function(params) {
              // build a color map as your need.
              var colorList = ['#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0', '#2297F0'];
              return colorList[params.dataIndex]
            },
            barBorderRadius:[50,50,0,0]
          }
        },
        data:$scope.myChart4.series.data//[99, 21, 10, 4, 12,99, 21, 10, 4, 12]
      }]
    };
    $scope.myChart4.Chart.setOption(option4);
  }
  $scope.htmlReady5 = function() {
    // 基于准备好的dom，初始化echarts实例
    $scope.myChart5.Chart = echarts.init(document.getElementById('sbfx_kt5'));
//		// 指定图表的配置项和数据
    var option5 = {
      fn: true,
      tooltip: {
        trigger: 'item'
      },
      toolbox: {
        show: false,
        feature: {dataView: {show: true,readOnly: false},
          restore: {show: true},
          saveAsImage: {show: true}
        }
      },
      calculable: true,
      grid: {
        borderWidth: 0,
        bottom:60,
        top:'20%',
        y: 80,
        y2: 60
      },
      xAxis: [{
        type: 'category',
        show: true,
        data: $scope.myChart5.xAxis.data,//['英语系', '法语系', '日语系', '西班牙语系'],
        axisLabel: {
          margin:20,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          }

        },
        axisLine: {
          show: false
        },
        axisTick:{
          show:false
        }
      }],
      yAxis: [{
        type: 'value',
        show: false
      }],
      series: [{
        name: '',
        type: 'bar',
        barWidth: 18,
        label:{
          normal:{
            show:true,
            position:'top',
            formatter:'{c}' + '%',
            textStyle:{
              color:'#7f8fa4',
              fontSize:"14"
            }
          }
        },
        itemStyle: {
          normal: {
            color: function(params) {
              // build a color map as your need.
              var colorList = ['#2297F0', '#2297F0', '#2297F0', '#2297F0'];
              return colorList[params.dataIndex]
            },
            barBorderRadius:[50,50,0,0]
          }
        },
        data:$scope.myChart5.series.data// [99, 21, 10, 4]
      }]
    };
    $scope.myChart5.Chart.setOption(option5);
  }
	Init_load();
	
	
	//   搜索条件   1
	$scope.searchData = {
		//  位置类型
		"SiteType": "",
		//  设备位置id：字符串
		"SiteId": "",
		//   设备型号id：字符串
		"ModelId": "",
		//   关键词
		"KeyWord": ""
	}
	var data_color = ['#2297F0', '#F5745B', '#FA4563', '#FFD400', '#40557D'];
	$scope.data_color = data_color;

	var legend_list = [{
		ModelName: "总台数"
	}, {
		ModelName: "停用中"
	}, {
		ModelName: "预警中"
	}, {
		ModelName: "故障中"
	}, {
		ModelName: "离线"
	}];
	//   柱格式
	var data_column = {
			'value': 20,
			'itemStyle': {
				'normal': {
					'color': '#fff'
				}
			}
		}
		//  option
	var options6 = {
		title: {
			show:false,
			text: '空调',
			textStyle: {
				color: '#fff',
				fontSize: 20
			}

		},
		grid:{
			left:30,
			right:'5%',
			top:'1%',
			bottom:30,
		},
		xAxis: {
	        type: 'value',
			data: [100,100,100,100,100],
			show:false,
			axisLabel: {
				inside: false,
				textStyle: {
					fontSize: 11,
					color: '#8DA1C1'
				}
			},
			axisTick: {
				show: false
			},
			axisLine: {
				show: false
			},
			z: 10
		},
		yAxis: {
			type: 'category',
			data: [100,100,100,100,100],
			show:false,
			axisLine: {
				show: false
			},
			axisTick: {
				show: false
			},
			axisLabel: {
				textStyle: {
					color: '#8DA1C1'
				}
			},
			min: 0,
			max: 'dataMax',
			splitLine: {
				show: true,
				lineStyle: {
					color: ['#ccc']
				}
			}
		},
		series: [{ // For shadow
			type: 'bar',
			itemStyle: {
				normal: {
					color: '#dfe2e5',
					barBorderRadius: 30
				}
			},
      		silent:true,
			barGap: '-100%',
			barCategoryGap: '40%',
			data: [100,100,100,100,100],
			animation: false,
			barWidth: '20%'
		}, {
			type: 'bar',
			barWidth: '20%',
			label: {
				normal: {
					show: true,
					position: [-20,-5],
					textStyle: {
						color: '#7f8fa4',
						fontSize: '12'
					}
				}
			},
			itemStyle: {
				normal: {
					color: '#ccc',
					borderWidth:11,
					barBorderRadius: 30
			  }
		  },
			data: [
  			{
				name:'总台数',
				value: 0,
				itemStyle: {
          		normal: {
					color: data_color[4]
            	}
          	}
			},
			{
				value: 0,
				itemStyle: {
					normal: {
						color: data_color[3]
            		}
          		}
			}, {
				value: 0,
				itemStyle: {
          normal: {
						color: data_color[2]
            }
          }
			}, {
				value: 0,
				itemStyle: {
          normal: {
						color: data_color[1]
            }
          }
			}, {
				value: 0,
				itemStyle: {
          normal: {
						color: data_color[0]
            }
          }
        }]
		}]
	};
	
	//   设备分析
	var getDeviceQty = function(sitetype, siteid, modelid) {
		var url = config.HttpUrl + "/device/getDeviceQty";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token
			},
			Para: {
				SiteType: sitetype,
				SiteId: siteid.toString(),
				ModelId: modelid.toString()
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.inDeviceItems = data.Result.Data;
				for(var i = 0; i < $scope.inDeviceItems.length; i++) {
					//    取最大值
					var max_v = 0;
					var max_v_arr = [];
					max_v_arr[0] = $scope.inDeviceItems[i].TotalQty;
					max_v_arr[1] = $scope.inDeviceItems[i].StopQty;
					max_v_arr[2] = $scope.inDeviceItems[i].AlertQty;
					max_v_arr[3] = $scope.inDeviceItems[i].FaultQty;
					max_v_arr[4] = $scope.inDeviceItems[i].OfflineQty;
					for(var b = 0; b < 5; b++){
						if(max_v_arr[b] > max_v){
							max_v = max_v_arr[b];
						}
					}
					//
					options6.series[1].name = $scope.inDeviceItems[i].ModelName;
					options6.series[1].data[0].value = $scope.inDeviceItems[i].TotalQty;
					options6.series[1].data[1].value = $scope.inDeviceItems[i].StopQty;
					options6.series[1].data[2].value = $scope.inDeviceItems[i].AlertQty;
					options6.series[1].data[3].value = $scope.inDeviceItems[i].FaultQty;
					options6.series[1].data[4].value = $scope.inDeviceItems[i].OfflineQty;
					//  放入最大值
					options6.series[0].data[0] = max_v;
					options6.series[0].data[1] = max_v;
					options6.series[0].data[2] = max_v;
					options6.series[0].data[3] = max_v;
					options6.series[0].data[4] = max_v;
					//   放入x标签
					options6.xAxis.data[0] = $scope.inDeviceItems[i].TotalQty;
					options6.xAxis.data[1] = $scope.inDeviceItems[i].StopQty;
					options6.xAxis.data[2] = $scope.inDeviceItems[i].AlertQty;
					options6.xAxis.data[3] = $scope.inDeviceItems[i].FaultQty;
					options6.xAxis.data[4] = $scope.inDeviceItems[i].OfflineQty;

					console.log(options6);
					$scope.inDeviceItems[i].option = JSON.stringify(options6);
				}
//				$scope.inDeviceItems = [{
//					option: ""
//				}];
				//$scope.inDeviceItems[0].option = JSON.stringify(optione);
      } else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
}
	$scope.navlist = [];
	$scope.geiappcode = '';
	var getappList = function (code) {
		var url = config.HttpUrl + "/getapp";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "PcWEB",
			"Token": config.GetUser().Token
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.navlist = data.Result;
				for(var i = 0; i < $scope.navlist.length; i++){
					switch($scope.navlist[i].Modulecode){
						case 'sbgl':
						$scope.geiappcode = $scope.navlist[i].Modulecode;
						break;
						case 'jcsj.yhgl.jsgl':
						$scope.geiappcode = $scope.navlist[i].Modulecode;
						break;
						case 'qxgl':
						$scope.geiappcode = $scope.navlist[i].Modulecode;
						break;
					}
				}
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}	
	//  run
	var run = function() {
		getDeviceQty($scope.searchData.SiteType, $scope.searchData.SiteId, $scope.searchData.ModelId);
		time();
		getappList($scope.geiappcode);
	}
	run();

}]);
