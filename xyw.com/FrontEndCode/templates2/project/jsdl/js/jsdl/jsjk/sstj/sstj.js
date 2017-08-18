'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

/*   实时统计     */
app.controller('jsdlJsjkSstjContr', ['$scope', '$rootScope', '$stateParams', 'httpService', '$modal', '$interval', '$location', '$filter', function($scope, $rootScope, $stateParams, httpService, $modal, $interval, $location, $filter) {
	console.log("实时统计:" + $location.search().ClassroomId);
	//$scope.Classroomid=$location.search().classroom;
	$scope.Begindatestr = "";
	$scope.Endatestr = "";
	$scope.YData = []; //定义Y轴数组
	$scope.XData = ['07:00', '08:00', '09:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00', '17:00', '18:00', '19:00', '20:00', '21:00', '22:00', '23:00']; //定义X轴数组
	$scope.seriesColor = ['#65EAD1', '#FF917B', '#FCFFC1', '#F469A9', '#E3E3E3', '#65EAD1', '#FF917B', '#FCFFC1', '#F469A9', '#E3E3E3'];
	$scope.series = []; //定义图表容器对象集合数组
	var Classroomnull = {
		Classroomid: $location.search().ClassroomId,
		Seatsnumbers: 0, //教室内的座位数
		Sumnumbers: 0, //教室内的人数
		Percentage: "", //教室内的入座率
		Classroomstate: 0, //教室的状态
		ClassroomstateStr: "" //教室的状态
	};
	$scope.ClassroomObject = {
		Classroom: Classroomnull,
		ClassroomData: []
	};
	var ToPercentage = function(number) {
		return(Math.round(number * 10000) / 100).toFixed(2) + '%';
	};
	//   教室导流 教室详情 教室信息查询
	var LoadClassroomInfo = function() {
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			"id": $scope.ClassroomObject.Classroom.Classroomid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.ClassroomObject.Classroom = data.Result;
				$scope.ClassroomObject.Classroom.Percentage = ToPercentage((data.Result.Sumnumbers / data.Result.Seatsnumbers));
				if(data.Result.Classroomstate == -1) {
					$scope.ClassroomObject.Classroom.ClassroomstateStr = "不可用";
				} else if(data.Result.Classroomstate == 1) {
					$scope.ClassroomObject.Classroom.ClassroomstateStr = "上课中";
				} else {
					$scope.ClassroomObject.Classroom.ClassroomstateStr = "开放中";
				}

				//   实时统计 cvs
				$scope.jsdl_sstj.options.series[1].data[0].value = $scope.ClassroomObject.Classroom.Sumnumbers;
				$scope.jsdl_sstj.options.series[1].data[1].value = $scope.ClassroomObject.Classroom.Seatsnumbers;
				$scope.jsdl_sstj.options = JSON.stringify($scope.jsdl_sstj.options);
			} else {
				$scope.ClassroomObject.Classroom = Classroomnull;
			}
		}, function(reason) {}, function(update) {});
	};
	//   教室导流 教室详情 教室内人员列表信息查询
	var LoadClassroomPeopleInfo = function() {
		VerifyValue();
		var url = config.HttpUrl + "/basicset/getclassroompeopleinfo";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Classroomid": Number($scope.ClassroomObject.Classroom.Classroomid),
			"Begindate": $scope.Begindatestr,
			"Enddate": $scope.Endatestr
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.ClassroomObject.ClassroomData = data.Result;

			} else {
				$scope.ClassroomObject.ClassroomData = [];
			}
		}, function(reason) {}, function(update) {});
	};
	$scope.Init_load = function() {
		LoadClassroomInfo();
		LoadClassroomPeopleInfo();
		LoadClassRoomChats();
	};
	$scope.loadValdata = function() {
		LoadClassroomPeopleInfo();
		LoadClassRoomChats();
	};
	var getNowFormatDate = function() {
		var date = new Date();
		var seperator1 = "-";
		var seperator2 = ":";
		var month = date.getMonth() + 1;
		var strDate = date.getDate();
		if(month >= 1 && month <= 9) {
			month = "0" + month;
		}
		if(strDate >= 0 && strDate <= 9) {
			strDate = "0" + strDate;
		}
		var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate;
		console.log(currentdate);
		return currentdate;
	};
	var VerifyValue = function() {
		if($scope.Begindatestr == "") {
			$scope.Begindatestr = getNowFormatDate() + " 00:00:00";
		}
		if($scope.Endatestr == "") {
			$scope.Endatestr = getNowFormatDate() + " 23:59:59";
		}
	};
	var LoadClassRoomChats = function() {

		/*---------------------开始获取图表数据---------------------------*/
		var url = config.HttpUrl + "/basicset/getclassroompeoplecount";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Classroomid": Number($scope.ClassroomObject.Classroom.Classroomid),
			"Begindate": $scope.Begindatestr,
			"Enddate": $scope.Endatestr
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("开始获取图表数据", data);
			if(data.Rcode == "1000") {
				ConvertDataChat(data.Result);
				//
				$scope.jsdl_qst.options.legend.data = $scope.YData;
				$scope.jsdl_qst.options.xAxis[0].data = $scope.XData;
				$scope.jsdl_qst.options.series = $scope.series;
				console.log($scope.jsdl_qst.options);
			} else {
				//
			}
		}, function(reason) {}, function(update) {});
		/*---------------------结束获取图表数据---------------------------*/
	}

	var ConvertDataChat = function(d) {
		var datalable = []; //定义使用时间分析数组
		var dataspan = []; //定义图表容器对象集合数组
		var Xdatalable = []; //定义X轴数组
		var Maxval = 7; //最大值
		var Minval = 23; //最小值
		/*--------------开始填充图表内的标识点----------------*/
		for(var k = 0; k < d.length; k++) {
			var temp = "";
			temp = d[k].Dateymd.substr(0, 4), temp += "-" + d[k].Dateymd.substr(4, 2), temp += "-" + d[k].Dateymd.substr(6, 2);
			d[k].Dateymd = temp;

			if(Number(d[k].Dateh) > Maxval) {
				Maxval = Number(d[k].Dateh);
			}
			if(Number(d[k].Dateh) < Minval) {
				Minval = Number(d[k].Dateh);
			}
			if(k == 0) {
				datalable.push(d[k].Dateymd);
			} else {
				var isadd = true;
				for(var p = 0; p < datalable.length; p++) {
					if(d[k].Dateymd == datalable[p]) {
						isadd = false;
						break;
					}
				}
				if(isadd) {
					datalable.push(d[k].Dateymd);
				}
			}
			if(Xdatalable.join("|").indexOf(d[k].Dateh) < 0) {
				Xdatalable.push(d[k].Dateh);
			}
		}
		console.log(Maxval);
		console.log(Minval);
		$scope.XData = Xdatalable;
		$scope.YData = datalable;
		/*--------------开始填充图表内的标识点----------------*/
		/*--------------开始填充图表内的数值----------------*/
		var objseries = {
			name: '',
			type: 'line',
			symbolSize:10,
			smooth:false,
			areaStyle: {normal: {opacity:0.4}},
			itemStyle: {
				normal: {
					areaStyle: {
						type: 'default'
					}
				}
			},
			data: [], //点位数据
			datah: [] //时间数据
		};
		$scope.XData = [];
		for(var p = 0; p < datalable.length; p++) {
			objseries = {
				name: '',
				type: 'line',
				symbolSize:10,
				smooth:false,
				areaStyle: {normal: {opacity:0.4}},
				itemStyle: {
					normal: {
						areaStyle: {
							type: 'default',
							'color': ''
						}
					}
				},
				data: [], //点位数据
				datah: [] //时间数据
			};
			objseries.itemStyle.normal.areaStyle.color = $scope.seriesColor[p];

			objseries.name = datalable[p];
			//objseries.name = temp;
			var ix = 0;
			for(var i = Minval; i <= Maxval; i++) {
				for(var k = 0; k < d.length; k++) {
					if(d[k].Dateymd == datalable[p]) {
						if(i == Minval) {
							if(i == Number(d[k].Dateh)) {
								objseries.data.push(d[k].Sumnumbers);
							} else {
								objseries.data.push(0);
							}
						} else {
							if(i == Number(d[k].Dateh)) {
								objseries.data[ix] = d[k].Sumnumbers;
							} else {
								if(objseries.data[ix] == 0 || objseries.data[ix] == '' || objseries.data[ix] == undefined) {
									objseries.data[ix] = 0;
								}
							}
						}
					}
				}
				if(p == 0) {
					var temp = "";
					i < 10 ? temp = "0" + i + ":00" : temp = i + ":00";
					$scope.XData.push(temp);
				}
				ix++;
			}
			dataspan.push(objseries);
		}
		console.log("dataspan", dataspan);
		$scope.series = dataspan;
		/*--------------结束填充图表内的数值----------------*/
	};
	//	var getNowFormatDate = function() {
	//      var date = new Date();
	//      var seperator1 = "-";
	//      var seperator2 = ":";
	//      var month = date.getMonth() + 1;
	//      var strDate = date.getDate();
	//      if (month >= 1 && month <= 9) {
	//          month = "0" + month;
	//      }
	//      if (strDate >= 0 && strDate <= 9) {
	//          strDate = "0" + strDate;
	//      }
	//      var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate
	//          + " " + date.getHours() + seperator2 + date.getMinutes()
	//          + seperator2 + date.getSeconds();
	//      return currentdate;
	//  };

	$scope.htmlReady = function() {
		// 基于准备好的dom，初始化echarts实例
		var myChart = echarts.init(document.getElementById('main'));
		// 指定图表的配置项和数据
		var option = {
			title: {
				text: '',
				subtext: ''
			},
			grid: {
				left: '2%',
				right: '2%',
				bottom: '3%',
				containLabel: true
			},
			tooltip: {
				trigger: 'axis',
				formatter: function(params) {
					if(!params[0]) {} else {
						var temp = "时间：" + params[0].name + "<br />";
					}
					for(var a in params) {
						temp += "<b style='display:inline-block;width: 8px;height: 8px;vertical-align: top;margin:6px 5px 0 0;background-color: " + params[a].color + ";border-radius: 50%;'></b>" + params[a].seriesName + "：" + params[a].value + "人<br />";
					}
					return temp;
				}
			},
			legend: {
				data: $scope.YData
			},
			toolbox: {
				show: true,
				feature: {
					mark: {
						show: true
					},
					//dataView : {show: true, readOnly: false},
					//magicType : {show: true, type: ['line', 'bar', 'stack', 'tiled']},
					//restore : {show: true},
					saveAsImage: {
						show: true
					}
				}
			},
			calculable: true,
			xAxis: [{
				type: 'category',
				boundaryGap: false,
				data: $scope.XData
			}],
			yAxis: [{
				type: 'value',
				minInterval: 1
			}],
			series: $scope.series
		};
		// 使用刚指定的配置项和数据显示图表。
		myChart.setOption(option);
	};

	//    可连续查询 天 数  //  天
	var dataLangthNumber = 15;

	var start = {
		dateCell: "#begindate",
		format: "YYYY-MM-DD hh:mm:ss",
		isTime: true,
		minDate: "2015-12-31",
		initAddVal: [2],
		festival: true,
		maxDate: jeDate.now(0),
		isinitVal: false,
		choosefun: function(elem, datas) {
			$scope.Begindatestr = datas;
			//   选择的开始时间加dataLangthNumber 天 时间戳
			var endTemp = Date.parse(new Date(datas)) + dataLangthNumber * 24 * 3600 * 1000;
			if(endTemp < Date.parse(new Date($scope.Endatestr)) && Date.parse(new Date(datas)) <= Date.parse(new Date($scope.Endatestr))) {
				$scope.Endatestr = $filter('date')(endTemp, 'yyyy-MM-dd HH:mm:ss');
				$("#enddate").val($scope.Endatestr);
			} else if(Date.parse(new Date(datas)) > Date.parse(new Date($scope.Endatestr))) {
				$scope.Endatestr = $filter('date')(Date.parse(new Date(datas)), 'yyyy-MM-dd') + " 23:59:59";
				$("#enddate").val($scope.Endatestr);
			}
			//  查询
			$scope.loadValdata();
		},
		okfun: function(elem, datas) {
			$scope.Begindatestr = datas;
			//   选择的开始时间加dataLangthNumber 天 时间戳
			var endTemp = Date.parse(new Date(datas)) + dataLangthNumber * 24 * 3600 * 1000;
			if(endTemp < Date.parse(new Date($scope.Endatestr)) && Date.parse(new Date(datas)) <= Date.parse(new Date($scope.Endatestr))) {
				$scope.Endatestr = $filter('date')(endTemp, 'yyyy-MM-dd HH:mm:ss');
				$("#enddate").val($scope.Endatestr);
			} else if(Date.parse(new Date(datas)) > Date.parse(new Date($scope.Endatestr))) {
				$scope.Endatestr = $filter('date')(Date.parse(new Date(datas)), 'yyyy-MM-dd') + " 23:59:59";
				$("#enddate").val($scope.Endatestr);
			}
			//  查询
			$scope.loadValdata();
		},
		clearfun: function(elem, datas) {
			$scope.Begindatestr = "";
		}
	};

	var end = {
		dateCell: "#enddate",
		format: "YYYY-MM-DD hh:mm:ss",
		isTime: true,
		minDate: "2015-12-31",
		maxDate: jeDate.now(0),
		isinitVal: false,
		choosefun: function(elem, datas) {
			$scope.Endatestr = datas;
			//   选择的开始时间加dataLangthNumber 天 时间戳
			var endTemp = Date.parse(new Date(datas)) - dataLangthNumber * 24 * 3600 * 1000;
			if(endTemp > Date.parse(new Date($scope.Begindatestr)) && Date.parse(new Date(datas)) >= Date.parse(new Date($scope.Begindatestr))) {
				$scope.Begindatestr = $filter('date')(endTemp, 'yyyy-MM-dd HH:mm:ss');
				$("#begindate").val($scope.Begindatestr);
			} else if(Date.parse(new Date(datas)) < Date.parse(new Date($scope.Begindatestr))) {
				$scope.Begindatestr = $filter('date')(Date.parse(new Date(datas)), 'yyyy-MM-dd') + " 23:59:59";
				$("#begindate").val($scope.Begindatestr);
			}
			//  查询
			$scope.loadValdata();
			//start.maxDate = datas;
		},
		okfun: function(elem, datas) {
			$scope.Endatestr = datas;

			//   选择的开始时间加dataLangthNumber 天 时间戳
			var endTemp = Date.parse(new Date(datas)) - dataLangthNumber * 24 * 3600 * 1000;
			if(endTemp > Date.parse(new Date($scope.Begindatestr)) && Date.parse(new Date(datas)) >= Date.parse(new Date($scope.Begindatestr))) {
				$scope.Begindatestr = $filter('date')(endTemp, 'yyyy-MM-dd HH:mm:ss');
				$("#begindate").val($scope.Begindatestr);
			} else if(Date.parse(new Date(datas)) < Date.parse(new Date($scope.Begindatestr))) {
				$scope.Begindatestr = $filter('date')(Date.parse(new Date(datas)), 'yyyy-MM-dd') + " 00:00:00";
				$("#begindate").val($scope.Begindatestr);
			}
			//  查询
			$scope.loadValdata();
		},
		clearfun: function(elem, datas) {
			$scope.Endatestr = "";
		}
	};

	//   选择时间
	$scope.beginDate = function() {
		jeDate(start);

	}
	$scope.endDate = function() {
		jeDate(end);
	}

	$scope.$watch('ClassroomObject.ClassroomData', function() {
		setTimeout(function() {
			$('#jsjk_sstj').trigger('footable_redraw');
		}, 1000);

	});

	//  -------------------------  实时统计 cvs  --------------------------
	//   echarts
	$scope.jsdl_sstj = {
		"options": {}
	};
	//
	$scope.jsdl_sstj.options = {
		color: ['#40557D', '#289DF5'],
		grid: {
			top: 0,
			bottom: 0,
			left: 0,
			right: 0
		},
		series: [{
			//name:'总座位数',
			type: 'pie',
			radius: ['85%', '100%'],
			avoidLabelOverlap: true,
			hoverAnimation: false,
			animation: false,
			labelLine: {
				normal: {
					show: false
				}
			},
			data: [{
				value: 1
			}]
		}, {
			//name:'已入座',
			type: 'pie',
			radius: ['85%', '100%'],
			avoidLabelOverlap: true,
			hoverAnimation: false,
			labelLine: {
				normal: {
					show: false
				}
			},
			data: [{
				value: 0
			}, {
				value: 100
			}]
		}]
	};

	//  -------------------------  /实时统计 cvs  --------------------------

	//  ------------------------  趋势图  --------------------------------
	//   echarts
	$scope.jsdl_qst = {
		"options": {}
	};
	$scope.jsdl_qst.options = {
		title: {
			//text: '趋势图（单位：人）'
		},
		color: ['#00AAFF', '#886CE6'],
		tooltip: {
			trigger: 'axis',
			formatter: function(params) {
				if(!params[0]) {} else {
					var temp = "时间：" + params[0].name + "<br />";
				}
				for(var a in params) {
					temp += "<b style='display:inline-block;width: 8px;height: 8px;vertical-align: top;margin:6px 5px 0 0;background-color: " + params[a].color + ";border-radius: 50%;'></b>" + params[a].seriesName + "：" + params[a].value + "人<br />";
				}
				return temp;
			}
		},
		legend: {
			data: []
		},
		grid: {
			left: '3%',
			right: '4%',
			bottom: '3%',
			containLabel: true
		},
		axisPointer:{
			show:false
		},
		xAxis: [{
			type: 'category',
			boundaryGap: false,
			axisLine: {
				show: false,
				lineStyle: {
					color: '#7F8FA4'
				}
			},
			axisTick: {
				show: false
			},
			data: []
		}],
		yAxis: [{
			axisLine: {
				show: false,
				lineStyle: {
					color: '#7F8FA4'
				}
			},
			axisTick: {
				show: false
			},
			splitLine: {
				lineStyle: {
          color: '#7F8FA4',
					opacity: 0.2
				}
			},
			nameTextStyle: {
				color: '#7F8FA4'
			},
			min: 0,
			minInterval:1,
			type: 'value'
		}],
		series: [

			//	        {
			//	            name:'直接访问',
			//	            type:'line',
			//	            stack: '总量',
			//				symbolSize:10,
			//				smooth:false,
			//	            itemStyle:{
			//	                normal:{
			//	                    width:10
			//	                }
			//	            },
			//	            markLine:{
			//	                silent:false
			//	            },
			//	            areaStyle: {normal: {}},
			//	            data:[320, 332, 301, 334, 390, 330, 320]
			//	        },
			//	        {
			//	            name:'搜索引擎',
			//	            type:'line',
			//	            stack: '总量',
			//				symbolSize:10,
			//				smooth:false,
			//	            areaStyle: {normal: {}},
			//	            data:[820, 932, 901, 934, 1290, 1330, 1320]
			//	        }
		]
	};
	//  ------------------------  /趋势图  --------------------------------

	VerifyValue();
	$scope.Init_load();
}]);
