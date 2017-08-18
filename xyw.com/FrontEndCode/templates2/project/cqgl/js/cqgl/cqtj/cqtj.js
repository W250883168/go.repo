'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 出勤统计
 */
app.filter('propsFilter', function() {
    return function(items, props) {
        var out = [];
        if (angular.isArray(items)) {
            items.forEach(function(item) {
                var itemMatches = false;
                var keys = Object.keys(props);
                for (var i = 0; i < keys.length; i++) {
                    var prop = keys[i];
                    try{
                        if(typeof(props[prop])=="number" ){
                            var text = props[prop];
                            if (item[prop].toString()==text.toString() ) {
                                itemMatches = true;
                                break;
                            }
                        }else{
                            var text = props[prop].toString().toLowerCase();
                            if (item[prop].toString().toLowerCase().indexOf(text) !== -1) {
                                itemMatches = true;
                                break;
                            }
                        }
                    }
                    catch(e){}
                }
                if (itemMatches) {
                    out.push(item);
                }
            });
        } else {
            out = items;
        }
        return out;
    };
});
/*   出勤统计      */
app.controller('cqglCqtjContr', ['$scope', 'httpService','toaster', function($scope, httpService,toaster) {
	console.log("出勤统计")
	$scope.selectdata=[{name:"昨天",value:1},{name:"近3天",value:3},{name:"近7天",value:7},{name:"近1个月",value:30},{name:"近3个月",value:90},{name:"近6个月",value:180},{name:"近1年",value:365}];
	$scope.gradelist=[{name:"大一",value:0},{name:"大二",value:1},{name:"大三",value:2},{name:"大四",value:3}];
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

    //$scope.changeselect=function(item,index){
    //console.log('item',item);
    //    switch(index){
    //        case 0:
    //          $scope.myChart1.MajorItem={};
    //          $scope.myChart4.MajorItem = item.Id
    //          $scope.myChart5.MajorItem = item.Id
    //        break;
    //        case 1:
    //          $scope.myChart4.CollegeItem = item.Id
    //          $scope.myChart5.CollegeItem = item.Id
    //          break;
    //        default:
    //          break;
    //    }
    //};

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
				data: $scope.myChart1.xAxis.data//[99, 21, 10, 4, 12]
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
}]);
