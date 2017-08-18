'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

/*   设备分析     */
app.controller('sbglSbfxContr', ['$scope', '$modal', 'httpService', '$filter', '$rootScope', '$location','toaster', function($scope, $modal, httpService, $filter, $rootScope, $location,toaster) {
	console.log("设备分析")

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
	//   按设备分类
	$scope.searchData2 = {
		//   开始时间
			"fromTime": "",
		//   结束时间
			"toTime": "",
		//  位置类型
			"SiteType": "",
		//  设备位置id：字符串
			"SiteId": "",
		//   设备型号id：字符串
			"ModelId": "",
		//   使用时间-日平均使用时间、总使用时间
			"oneAll": 2,
			// ----  使用时间-日平均使用时间、总使用时间
			"oneAllItem": {"val":2,"title":"总使用时间"},
			// --- - 使用时间-日平均使用时间、总使用时间
			"oneAllItems": [
				{"val":1,"title":"日平均使用时间"},
				{"val":2,"title":"总使用时间"}
			],
		//   使用时间-日平均使用时间、总使用时间- 天数
			"oneAllDay": 1,
		//   最近几个月使用时间
			"howMonth": 1,
			// ----  最近几个月使用时间
			"howMonthItem": {"val":1,"title":"近1个月"},
			// ----  最近几个月使用时间
			"howMonthItems": [
				{"val":1,"title":"近1个月"},
				{"val":3,"title":"近3个月"},
				{"val":6,"title":"近6个月"},
				{"val":12,"title":"近1年"}
			]
	}
	//   按设备位置
	$scope.searchData3 = {
		//   开始时间
		"fromTime": "",
		//   结束时间
		"toTime": "",
		//  位置类型
		"SiteType": "",
		//  设备位置id：字符串
		"SiteId": "",
		//   设备型号id：字符串
		"ModelId": "",
		//   使用时间-日平均使用时间、总使用时间
		"oneAll": 2,
		// ----  使用时间-日平均使用时间、总使用时间
		"oneAllItem": {"val":2,"title":"总使用时间"},
		// --- - 使用时间-日平均使用时间、总使用时间
		"oneAllItems": [
			{"val":1,"title":"日平均使用时间"},
			{"val":2,"title":"总使用时间"}
		],
		//   使用时间-日平均使用时间、总使用时间- 天数
		"oneAllDay": 1,
		//   最近几个月使用时间
		"howMonth": 1,
		// ----  最近几个月使用时间
		"howMonthItem": {"val":1,"title":"近1个月"},
		// ----  最近几个月使用时间
		"howMonthItems": [
			{"val":1,"title":"近1个月"},
			{"val":3,"title":"近3个月"},
			{"val":6,"title":"近6个月"},
			{"val":12,"title":"近1年"}
		]
	}

	//   返回 位置数据
	$scope.backAdd = {};
	//   位置HTML显示
	$scope.addAdd = "";
	//   设备HTML显示
	$scope.deviceText = "";

	//   返回 位置数据
	$scope.backAdd2 = {};
	//   位置HTML显示
	$scope.addAdd2 = "";
	//   设备HTML显示
	$scope.deviceText2 = "";

	//   返回 位置数据
	$scope.backAdd3 = {};
	//   位置HTML显示
	$scope.addAdd3 = "";
	//   设备HTML显示
	$scope.deviceText3 = "";

	//   当前设备统计
	$scope.inDeviceItems = [];
	//   按设备分类
	$scope.inDeviceItems2 = {};
	//   按设备位置
	$scope.inDeviceItems3 = [];

	//   formatter 按设备分类 与 按设备位置 数字转小时
	//   options：options echarts
	$scope.formatter = function(options, element) {
		if(element[0].id == "typeSite" || element[0].id == "typeTime") {
			//    修改柱上文字
			options.series[0].label.normal.formatter = function(params) {
					var temp_1 = (params.value / 3600);
					if(temp_1 >= 0){
						if(temp_1 % 1 > 0){
							return temp_1.toFixed(2);
						}else{
							return temp_1;
			}
					}else{
						return 0;
					}
				}
			//    修改鼠标上文字
			options.tooltip.formatter = function(params, b, c) {
				var temp_1 = (params.value / 3600);
				if(temp_1 >= 0){
					if(temp_1 % 1 > 0){
						return params.name + "：" + temp_1.toFixed(2) + "小时";
					}else{
						return params.name + "：" + temp_1 + "小时";
			}
				}else{
					return 0 + "小时";
		}
			}
		}
    return options;
	}


	//   颜色
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
	var options = {
		title: {
			show:false,
			text: '空调',
			textStyle: {
				color: '#fff',
				fontSize: 20
			}

		},
		grid:{
			left:'5%',
			right:'5%',
			top:'1%',
			bottom:30,
		},
		xAxis: {
			data: [100,100,100,100,100],
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
					barBorderRadius: [30, 30, 0, 0]
				}
			},
      silent:true,
			barGap: '-100%',
			barCategoryGap: '40%',
			data: [100,100,100,100,100],
			animation: false,
			barWidth: '30%'
		}, {
			type: 'bar',
			barWidth: '30%',
			itemStyle: {
				normal: {
					color: '#ccc',
					barBorderRadius: [30, 30, 0, 0]
			  }
		  },
			data: [
      {
				name:'总台数',
				value: 0,
				itemStyle: {
          normal: {
						color: data_color[0]
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
						color: data_color[2]
            }
          }
			}, {
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
						color: data_color[4]
            }
          }
        }]
		}]

	};
	//console.log(options)
	var options2 = {
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
			bottom:40,
			top:'20%',
			y: 80,
			y2: 60
		},
		xAxis: [{
			type: 'category',
			show: true,
			data: ['总台数', '停用中', '预警中', '故障中', '离线'],
			axisLabel: {
				margin: 10,
				textStyle: {
					color: '#7f8fa4',
					fontSize: 14,
					fontWeight: ''
				}
			},
			axisLine: {
				show: false
			},
			axisTick: {
				show: false
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
			label: {
				normal: {
					show: true,
					position: 'top',
					formatter: function(a) {
						//console.log(a)
					},
					textStyle: {
						color: '#7f8fa4',
						fontSize: '14'
					}
				}
			},
			itemStyle:{
				normal:{
					barBorderRadius:[50,50,0,0]
						}
				},
			data: [{
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#65EAD1'
						}
					}
			}, {
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#FF917B'
						}
					}
			}, {
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#FCFFC1'
						}
					}
			}, {
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#F469A9'
						}
					}
			}, {
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#E3E3E3'
						}
					}
			}, {
				'value': 0,
				'itemStyle': {
					'normal': {
						'color': '#E3E3E3'
				}
				}
		}]
		}]
	};

	//   取前几个月
	//    number   月数
	var getGapMonth = function(number) {
		(Number(number) >= 0 && Number(number) <= 120) ? number = Number(number): number = 1;
		var now   = new Date();
		//   结束时间
		var month = now.getMonth();
		var day = now.getDate();
		(month + 1) < 10 ? month = "0" + (month + 1) : month = (month + 1);
		(day) < 10 ? day = "0" + (day) : day = (day);
		var end = now.getFullYear() + '-' + month + '-' + day + ' 23:59:59';
		//   开始时间
		now.setMonth(now.getMonth() - number);
		var month = now.getMonth();
		var day = now.getDate();
		(month + 1) < 10 ? month = "0" + (month + 1) : month = (month + 1);
		(day) < 10 ? day = "0" + (day) : day = (day);
		var be = now.getFullYear() + '-' + month + '-' + day + ' 00:00:00';
		//   天数
		//var DayNumber = Math.floor((Date.parse(new Date(end.replace(/\-/g, "/"))) - Date.parse(new Date(be.replace(/\-/g, "/")))) / 1000 / 86400);
		var DayNumber = number * 30;

		return {
			"fromTime": be,
			"toTime": end,
			"dayNumber": DayNumber
		};
	}

	//

	//    打开弹窗  选择学校 位置
	$scope.modalOpenClassroom = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			//console.log(selectedItem)
			if(!selectedItem) {
				$scope.addAdd = "";
				$scope.searchData.SiteType = "";
				$scope.searchData.SiteId = "";
			} else {
				$scope.backAdd = selectedItem;
				$scope.searchData.SiteType = $scope.backAdd.addCode;
				$scope.searchData.SiteId = $scope.backAdd.addId;
				$scope.addAdd = $scope.backAdd.add;
			}
			//  查询
			getDeviceQty($scope.searchData.SiteType, $scope.searchData.SiteId, $scope.searchData.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择设备型号
	$scope.modalOpenDevice = function(obj) {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device.html',
			controller: 'modalGetDeviceCtrl',
			resolve: {
				items: function() {
					return obj;
				}
			}
		});

		modalInstance.result.then(function(deviceItem) {
			//console.log(deviceItem)
			if(!deviceItem) {
				$scope.deviceText = "";
				$scope.searchData.ModelId = "";
			} else {
				$scope.searchData.ModelId = deviceItem.Id;
				$scope.deviceText = deviceItem.Name;
			}
			//  查询
			getDeviceQty($scope.searchData.SiteType, $scope.searchData.SiteId, $scope.searchData.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//////////////////////  2  ////////////////////////
	//    打开弹窗  选择学校 位置2
	$scope.modalOpenClassroom2 = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			//console.log(selectedItem)
			if(!selectedItem) {
				$scope.addAdd2 = "";
				$scope.searchData2.SiteType = "";
				$scope.searchData2.SiteId = "";
			} else {
				$scope.backAdd2 = selectedItem;
				$scope.searchData2.SiteType = $scope.backAdd2.addCode;
				$scope.searchData2.SiteId = $scope.backAdd2.addId;
				$scope.addAdd2 = $scope.backAdd2.add;
			}
			//  查询
			getUseTimeByModel($scope.searchData2.fromTime, $scope.searchData2.toTime, $scope.searchData2.SiteType, $scope.searchData2.SiteId, $scope.searchData2.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择设备型号
	$scope.modalOpenDevice2 = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device.html',
			controller: 'modalGetDeviceCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(deviceItem) {
			//console.log(deviceItem)
			if(!deviceItem) {
				$scope.deviceText2 = "";
				$scope.searchData2.ModelId = "";
			} else {
				$scope.searchData2.ModelId = deviceItem.Id;
				$scope.deviceText2 = deviceItem.Name;
			}
			//  查询
				getUseTimeByModel($scope.searchData2.fromTime, $scope.searchData2.toTime, $scope.searchData2.SiteType, $scope.searchData2.SiteId, $scope.searchData2.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	//    打开弹窗  选择学校 位置
	$scope.modalOpenClassroom3 = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem) {
				$scope.addAdd3 = "";
				$scope.searchData3.SiteType = "";
				$scope.searchData3.SiteId = "";
			} else {
				$scope.backAdd3 = selectedItem;
				$scope.searchData3.SiteType = $scope.backAdd3.addCode;
				$scope.searchData3.SiteId = $scope.backAdd3.addId;
				$scope.addAdd3 = $scope.backAdd3.add;
				//   不能为教室
				if($scope.backAdd3.addCode == 'classroom') {
					//   为教室时换到楼层
					$scope.searchData3.SiteType = $scope.backAdd3.addItems.floor.addCode;
					$scope.searchData3.SiteId = $scope.backAdd3.addItems.floor.addId;
					$scope.addAdd3 = $scope.backAdd3.addItems.floor.add;
				}
			}
			//  查询
			getUseTimeBySite($scope.searchData3.fromTime, $scope.searchData3.toTime, $scope.searchData3.SiteType, $scope.searchData3.SiteId, $scope.searchData3.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择设备型号
	$scope.modalOpenDevice3 = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device.html',
			controller: 'modalGetDeviceCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(deviceItem) {
			//console.log(deviceItem)
			if(!deviceItem) {
				$scope.deviceText3 = "";
				$scope.searchData3.ModelId = "";
			} else {
				$scope.searchData3.ModelId = deviceItem.Id;
				$scope.deviceText3 = deviceItem.Name;
			}
			//  查询
			getUseTimeBySite($scope.searchData3.fromTime, $scope.searchData3.toTime, $scope.searchData3.SiteType, $scope.searchData3.SiteId, $scope.searchData3.ModelId);
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

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
					options.series[1].name = $scope.inDeviceItems[i].ModelName;
					options.series[1].data[0].value = $scope.inDeviceItems[i].TotalQty;
					options.series[1].data[1].value = $scope.inDeviceItems[i].StopQty;
					options.series[1].data[2].value = $scope.inDeviceItems[i].AlertQty;
					options.series[1].data[3].value = $scope.inDeviceItems[i].FaultQty;
					options.series[1].data[4].value = $scope.inDeviceItems[i].OfflineQty;
					//  放入最大值
					options.series[0].data[0] = max_v;
					options.series[0].data[1] = max_v;
					options.series[0].data[2] = max_v;
					options.series[0].data[3] = max_v;
					options.series[0].data[4] = max_v;
					//   放入x标签
					options.xAxis.data[0] = $scope.inDeviceItems[i].TotalQty;
					options.xAxis.data[1] = $scope.inDeviceItems[i].StopQty;
					options.xAxis.data[2] = $scope.inDeviceItems[i].AlertQty;
					options.xAxis.data[3] = $scope.inDeviceItems[i].FaultQty;
					options.xAxis.data[4] = $scope.inDeviceItems[i].OfflineQty;

					console.log(options);
					$scope.inDeviceItems[i].option = JSON.stringify(options);
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

	//   按设备分类
	var getUseTimeByModel = function(fromtime, totime, sitetype, siteid, modelid) {
		var url = config.HttpUrl + "/device/getUseTimeByModel";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token
			},
			Para: {
				FromTime: fromtime,
				ToTime: totime,
				SiteType: sitetype,
				SiteId: siteid.toString(),
				ModelId: modelid.toString()
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log('按设备分类')
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.inDeviceItems2 = data.Result.Data;
				//
				var option_temp = angular.copy(options2);
				var temp_data = [];
				for(var i = 0; i < $scope.inDeviceItems2.length; i++) {
					//option_temp.xAxis.Data
					option_temp.xAxis[0].data[i] = $scope.inDeviceItems2[i].ModelName;

					data_column.value = parseInt($scope.inDeviceItems2[i].UseTime / $scope.searchData2.oneAllDay);
					//data_column.itemStyle.normal.color = data_color[i];
					data_column.itemStyle.normal.color = data_color[0];
					temp_data.push(angular.copy(data_column));
				}
				option_temp.series[0].name = $scope.deviceText2;
				option_temp.series[0].data = temp_data;
				$scope.inDeviceItems2.option = option_temp;
				console.log('按设备分类',$scope.inDeviceItems2.option);
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   按设备位置
	var getUseTimeBySite = function(fromtime, totime, sitetype, siteid, modelid) {
		var url = config.HttpUrl + "/device/getUseTimeBySite";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token
			},
			Para: {
				FromTime: fromtime,
				ToTime: totime,
				SiteType: sitetype,
				SiteId: siteid.toString(),
				ModelId: modelid.toString()
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log('按设备位置')
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.inDeviceItems3 = data.Result.Data;
				//
				var option_temp3 = angular.copy(options2);
				var temp_data = [];
				for(var i = 0; i < $scope.inDeviceItems3.length; i++) {
					//option_temp.xAxis.Data
					option_temp3.xAxis[0].data[i] = $scope.inDeviceItems3[i].SiteName;

					data_column.value = parseInt($scope.inDeviceItems3[i].UseTime / $scope.searchData3.oneAllDay);
					data_column.name = $scope.inDeviceItems3[i].SiteName;
					//data_column.itemStyle.normal.color = data_color[i];
					data_column.itemStyle.normal.color = data_color[0];
					temp_data.push(angular.copy(data_column));
				}
				option_temp3.series[0].name = $scope.deviceText3;
				option_temp3.series[0].data = temp_data;
				$scope.inDeviceItems3.option = option_temp3;
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    按设备分类
	$scope.selectChange2 = function(time) {
		switch(time) {
			case "time":
				$scope.searchData2.oneAll = $scope.searchData2.oneAllItem.val;
				//
				if($scope.searchData2.oneAll == 1) {
					$scope.searchData2.oneAllDay = getGapMonth($scope.searchData2.howMonth).dayNumber;
					getUseTimeByModel($scope.searchData2.fromTime, $scope.searchData2.toTime, $scope.searchData2.SiteType, $scope.searchData2.SiteId, $scope.searchData2.ModelId);
				} else {
					$scope.searchData2.oneAllDay = 1;
					getUseTimeByModel($scope.searchData2.fromTime, $scope.searchData2.toTime, $scope.searchData2.SiteType, $scope.searchData2.SiteId, $scope.searchData2.ModelId);
				}
			break;
			case "month":
				$scope.searchData2.howMonth = $scope.searchData2.howMonthItem.val;
				//
				$scope.searchData2.fromTime = getGapMonth($scope.searchData2.howMonth).fromTime;
				$scope.searchData2.toTime = getGapMonth($scope.searchData2.howMonth).toTime;
				getUseTimeByModel($scope.searchData2.fromTime, $scope.searchData2.toTime, $scope.searchData2.SiteType, $scope.searchData2.SiteId, $scope.searchData2.ModelId);
			break;
		}

	}

	//    按设备分类
	$scope.selectChange3 = function(time) {
		switch(time) {
			case "time":
				$scope.searchData3.oneAll = $scope.searchData3.oneAllItem.val;
				//
				if($scope.searchData3.oneAll == 1) {
					$scope.searchData3.oneAllDay = getGapMonth($scope.searchData3.howMonth).dayNumber;
					getUseTimeBySite($scope.searchData3.fromTime, $scope.searchData3.toTime, $scope.searchData3.SiteType, $scope.searchData3.SiteId, $scope.searchData3.ModelId);
				} else {
					$scope.searchData3.oneAllDay = 1;
					getUseTimeBySite($scope.searchData3.fromTime, $scope.searchData3.toTime, $scope.searchData3.SiteType, $scope.searchData3.SiteId, $scope.searchData3.ModelId);
				}
			break;
			case "month":
				$scope.searchData3.howMonth = $scope.searchData3.howMonthItem.val;
				//
				$scope.searchData3.fromTime = getGapMonth($scope.searchData3.howMonth).fromTime;
				$scope.searchData3.toTime = getGapMonth($scope.searchData3.howMonth).toTime;
				getUseTimeBySite($scope.searchData3.fromTime, $scope.searchData3.toTime, $scope.searchData3.SiteType, $scope.searchData3.SiteId, $scope.searchData3.ModelId);
			break;
		}

	}

	//  run
	var run = function() {
		$scope.inDeviceItems2.option = options2;
		$scope.inDeviceItems3.option = options2;

		getDeviceQty($scope.searchData.SiteType, $scope.searchData.SiteId, $scope.searchData.ModelId);

		$scope.selectChange2('month');
		$scope.selectChange3('month');
	}
	run();

}]);
