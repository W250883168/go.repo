'use strict';
/**
 * Created by Administrator on 2016/9/8.
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
/*   导流分析     */
app.controller('jsdlDlfxContr', ['$scope', 'httpService', '$modal','$interval','toaster', function($scope, httpService, $modal,$interval,toaster) {
	//console.log("导流分析")

	$scope.Begindatestr="";
    $scope.Enddatestr="";
	$scope.Campusid=0;
	$scope.Buildingid=0;
	$scope.Floorsid=0;
	$scope.Campus={};
    $scope.Campus.Campuslist=[];
    $scope.Campus.selectCampusitem={};
    $scope.building={};
    $scope.building.buildinglist=[];
    $scope.building.selectbuildingitem={};
    $scope.floors={};
    $scope.floors.floorslist=[];
    $scope.floors.selectfloorsitem={};

	$scope.collegedata=[];//学院饼图统计
	$scope.majordata=[];//专业饼图统计
	$scope.classesdata=[];//班级饼图统计
	$scope.sexdata=[];//性别饼图统计
  $scope.inDeviceItems = [
    {
      title:"学院",
      dlfx:"dlfx_rq3",
      legend:['运达校区','测试校区555','ceshi321'],
      seriesData:[
        {
          value:75,
          name:'运达校区',
          itemStyle: {
            normal: {
              color: '#6ACF67'
            }
          }
        },
        {
          value:15,
          name:'测试校区555',
          itemStyle: {
            normal: {
              color: '#2297F0'
            }
          }
        },
        {
          value:10,
          name:'ceshi321',
          itemStyle: {
            normal: {
              color: '#5857CD'
            }
          }
        }
      ]
    },
    {
      title:"专业",
      dlfx:"dlfx_rq4",
      legend:['软件工程','计算机系列','英语系','数学系'],
      seriesData:[
        {value:65,name:'软件工程',
          itemStyle: {
            normal: {
              color: '#6ACF67'
            }
          }
        },
        {value:15,name:'计算机系列',
          itemStyle: {
            normal: {
              color: '#2297F0'
            }
          }
        },
        {value:5,name:'英语系',
          itemStyle: {
            normal: {
              color: '#5857CD'
            }
          }
        ,},
        {value:15,name:'数学系',
          itemStyle: {
            normal: {
              color: '#F469A9'
            }
          }
        }
      ]
    },
    {
      title:"班级",
      dlfx:"dlfx_rq5",
      legend:['测试班级','计算机班级','全能班','1801班'],
      seriesData:[
        {value:12,name:'测试班级',
          itemStyle: {
            normal: {
              color: '#6ACF67'
            }
          }
        },
        {value:48,name:'计算机班级',
          itemStyle: {
            normal: {
              color: '#2297F0'
            }
          }
        },
        {value:25,name:'全能班',
          itemStyle: {
            normal: {
              color: '#5857CD'
            }
          }
        },
        {value:15,name:'1801班',
          itemStyle: {
            normal: {
              color: '#F469A9'
            }
          }
        }
      ]
    },
    {
      title:"性别",
      dlfx:"dlfx_rq6",
      legend:['男生','女生'],
      seriesData:[
        {value:65,name:'男生',
          itemStyle: {
            normal: {
              color: '#6ACF67'
            }
          }
        },
        {value:35,name:'女生',
          itemStyle: {
            normal: {
              color: '#2297F0'
            }
          }
        }
      ]
    }
  ];//学院饼图统计
	$scope.StreamPeopledata = [{name:'',type:'bar',barWidth: '18',data:[],
		label: {
			normal: {
				show: true,
				position: 'top',
				formatter:'{c}',
				textStyle: {
					color: '#7f8fa4',
					fontSize: '14'
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
    }
	}];//人流柱状统计
	$scope.xAxisdata={};
	$scope.xAxisdata.data=[];//人流统计X轴数据
	$scope.xAxisdata.config=[];//人流统计X轴数据

	$scope.Xheatmap=[];//热力图X数据
	$scope.Yheatmap=[];//热力图Y数据

	$scope.changeselect=function(index){
        switch(index){
            case 0:$scope.building.selectbuildingitem={};//清除选中的楼栋
            case 1:$scope.floors.selectfloorsitem={};break;
            default:
                $scope.Campus.selectCampusitem={};
                $scope.building.selectbuildingitem={};
                $scope.floors.selectfloorsitem={};
        }
    };
	var Init_load=function()//初始化加载相关设置数据
	{
		var url = config.HttpUrl+"/basicset/getall";
        var data={};
        var promise =httpService.ajaxGet(url,null);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
                $scope.Campus.Campuslist=data.Result[0];
                $scope.building.buildinglist=data.Result[1];
                $scope.floors.floorslist=data.Result[2];
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
//      Load_value();
//		Load_StreamPeoplesAnalysis();
		$scope.queryfuncbut();
	};
	$scope.queryfuncbut=function(){
		Load_value();
        $scope.collegedata=[];//学院饼图统计
		$scope.majordata=[];//专业饼图统计
		$scope.classesdata=[];//班级饼图统计
		$scope.sexdata=[];//性别饼图统计
    $scope.StreamPeopledata = [{name:'',type:'bar',barWidth: '18',data:[18,63,20,10],
      label: {
        normal: {
          show: true,
          position: 'top',
          formatter:'{c}',
          textStyle: {
            color: '#7f8fa4',
            fontSize: '14'
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
      }
    }];
		$scope.xAxisdata={};
		$scope.xAxisdata.data=[];//人流统计X轴数据
		$scope.xAxisdata.config=[];//人流统计X轴数据
		Load_StreamPeoplesAnalysis();
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
		//console.log(currentdate);
		return currentdate;
	};
	var Load_value=function(){
		if($scope.Begindatestr==""){
			$scope.Begindatestr=getNowFormatDate()+" 00:00:00";
        }else{
        	$scope.Begindatestr=$("#begindate").val();
        }
        if($scope.Enddatestr==""){
        	$scope.Enddatestr=getNowFormatDate()+" 23:59:59";
        }else{
        	$scope.Enddatestr=$("#enddate").val();
        }
        if($scope.Campus.selectCampusitem !=undefined && $scope.Campus.selectCampusitem.Campusid !=null && $scope.Campus.selectCampusitem.Campusid !=""){
        	$scope.Campusid=Number($scope.Campus.selectCampusitem.Campusid);
        }
        if($scope.building.selectbuildingitem !=undefined && $scope.building.selectbuildingitem.Buildingid !=null && $scope.building.selectbuildingitem.Buildingid !=""){
        	$scope.Buildingid=Number($scope.building.selectbuildingitem.Buildingid);
        }
        if($scope.floors.selectfloorsitem !=undefined && $scope.floors.selectfloorsitem.Floorsid !=null && $scope.floors.selectfloorsitem.Floorsid !=""){
        	$scope.Floorsid=Number($scope.floors.selectfloorsitem.Floorsid);
        }
	};
	var Load_StreamPeoplesAnalysis=function()//加载统计数据
	{
		var url = config.HttpUrl+"/basicset/getstreampeoplesanalysis";
        var data={
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Os": "WEB",
			"Campusid": Number($scope.Campusid),
			"Buildingid": Number($scope.Buildingid),
			"Floorsid": Number($scope.Floorsid),
			"Begindate":$scope.Begindatestr,
			"Enddate":$scope.Enddatestr
		};
        var promise =httpService.ajaxPost(url,data);
        promise.then(function (data) {
            if(data.Rcode=="1000"){
				if(data.Result != null)DataBind(data.Result);
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
	};
	var DataBind=function(d)
	{
		//console.log(d);

			for(var c=0;c<d[1].length;c++){//统计学院数据
				$scope.collegedata.push({value:d[1][c].Valcount,name:d[1][c].Valname});
			}
			for(var m=0;m<d[2].length;m++){//统计专业数据
				$scope.majordata.push({value:d[2][m].Valcount,name:d[2][m].Valname});
			}
			for(var cs=0;cs<d[3].length;cs++){//统计班级数据
				$scope.classesdata.push({value:d[3][cs].Valcount,name:d[3][cs].Valname});
			}
			for(var s=0;s<d[4].length;s++){//统计性别数据
				$scope.sexdata.push({value:d[4][s].Valcount,name:GetSexName(d[4][s].Valname)});
			}
			for(var x=0;x<d[0].length;x++){//人流统计X轴数据
				$scope.xAxisdata.data.push(d[0][x].Valname);
				$scope.StreamPeopledata.data.push(d[0][x].Valcount);
			}

//		for(var i=0;i<d.length;i++){
//			var iscollegedataadd=true;
//			for(var c=0;c<$scope.collegedata.length;c++){//统计学院数据
//				if($scope.collegedata[c].name == d[i].Collegename && d[i].Collegename!='' && $scope.collegedata[c].name!=''){
//					$scope.collegedata[c].value=$scope.collegedata[c].value+d[i].Valcount;
//					iscollegedataadd=false;
//					break
//				}
//			}
//			var ismajordataadd=true;
//			for(var m=0;m<$scope.majordata.length;m++){//统计专业数据
//				if($scope.majordata[m].name==d[i].Majorname && d[i].Majorname!='' && $scope.collegedata[c].name!=''){
//					$scope.majordata[m].value=$scope.majordata[m].value+d[i].Valcount;
//					ismajordataadd=false;
//					break
//				}
//			}
//			var isclassesdataadd=true;
//			for(var cs=0;cs<$scope.classesdata.length;cs++){//统计班级数据
//				if($scope.classesdata[cs].name==d[i].Classesname && d[i].Classesname!='' && $scope.collegedata[c].name!=''){
//					$scope.classesdata[cs].value=$scope.classesdata[cs].value+d[i].Valcount;
//					isclassesdataadd=false;
//					break
//				}
//			}
//			var issexdataadd=true;
//			for(var s=0;s<$scope.sexdata.length;s++){//统计性别数据
//				if($scope.sexdata[s].name==GetSexName(d[i].Usersex) && $scope.sexdata[c].name!=''){
//					$scope.sexdata[s].value=$scope.sexdata[s].value+d[i].Valcount;
//					issexdataadd=false;
//					break
//				}
//			}
//			var isxAxisdataadd=true;
//			for(var x=0;x<$scope.xAxisdata.data.length;x++){//人流统计X轴数据
//				if(($scope.xAxisdata.data[x]==d[i].Floorname)||($scope.xAxisdata.data[x]==d[i].Buildingname)||($scope.xAxisdata.data[x]==d[i].Campusname)){
//						$scope.xAxisdata.config[x].value=$scope.xAxisdata.config[x].value+d[i].Valcount;
//						isxAxisdataadd=false;
//						break
//				}
//			}
//			if($scope.xAxisdata.data.length==0 || isxAxisdataadd){
//				if(Number($scope.Buildingid)>0){
//					$scope.xAxisdata.data.push(d[i].Floorname);
//					$scope.xAxisdata.config.push({name:d[i].Floorname,value:d[i].Valcount});
//				}else if(Number($scope.Campusid)>0){
//					$scope.xAxisdata.data.push(d[i].Buildingname);
//					$scope.xAxisdata.config.push({name:d[i].Buildingname,value:d[i].Valcount});
//				}else{
//					$scope.xAxisdata.data.push(d[i].Campusname);
//					$scope.xAxisdata.config.push({name:d[i].Campusname,value:d[i].Valcount});
//				}
//			}
//			if($scope.collegedata.length==0 || iscollegedataadd){
//				if(d[i].Collegename!=''){
//				$scope.collegedata.push({value:d[i].Valcount,name:d[i].Collegename});
//				}
//			}
//			if($scope.majordata.length==0 || ismajordataadd){
//				if(d[i].Majorname!=''){
//				$scope.majordata.push({value:d[i].Valcount,name:d[i].Majorname});
//				}
//			}
//			if($scope.classesdata.length==0 || isclassesdataadd){
//				if(d[i].Classesname!=''){
//				$scope.classesdata.push({value:d[i].Valcount,name:d[i].Classesname});
//				}
//			}
//			if($scope.sexdata.length==0 || issexdataadd){
//				if(GetSexName(d[i].Usersex)!=''){
//				$scope.sexdata.push({value:d[i].Valcount,name:GetSexName(d[i].Usersex)});
//				}
//			}
//		}
//		for(var k=0;k<$scope.xAxisdata.config.length;k++){
//			$scope.StreamPeopledata.data.push($scope.xAxisdata.config[k].value);
//		}

		$scope.htmlReady();
		//console.log($scope.StreamPeopledata);
	};
	var GetSexName = function(sex){
		var sexstr='保密';
		switch(sex){
			case "1":sexstr='女'; break;
			case "2":sexstr='男';break;
		}
		return sexstr;
	};
		//   HTML ready
	$scope.showFromDate = function() {
		jeDate({
			dateCell: "#begindate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.Begindatestr = val;
			},
			okfun: function(elem,val) {
				$scope.Begindatestr = val;
			},
			clearfun:function(elem, val) {
				$scope.Begindatestr = "";
			}
		});
	}

	$scope.showToDate = function() {
		jeDate({
			dateCell: "#enddate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.Enddatestr = val;
			},
			okfun: function(elem,val) {
				$scope.Enddatestr = val;
			},
			clearfun:function(elem, val) {
				$scope.Enddatestr = "";
			}
		});
	}

    $scope.htmlReady = function(){
	     // 基于准备好的dom，初始化echarts实例
        //var myChart1 = echarts.init(document.getElementById('dlfx_rq1'));
        var myChart2 = echarts.init(document.getElementById('dlfx_rq2'));
		    var myChart3 = echarts.init(document.getElementById('dlfx_rq3'));
		    var myChart4 = echarts.init(document.getElementById('dlfx_rq4'));
		    var myChart5 = echarts.init(document.getElementById('dlfx_rq5'));
		    var myChart6 = echarts.init(document.getElementById('dlfx_rq6'));
		var option2 = {
		    //color: ['#3398DB'],//控制数据的颜色
		    tooltip : {
		        trigger: 'item'
		    },
			//legend:{data:['直接访问']},
        grid: {
          borderWidth: 0,
          bottom:60,
          top:'20%',
          y: 80,
          y2: 60
        },
		    xAxis : [{
          type : 'category',
          data : ['A栋','B栋','C栋','D栋'],
		      axisLabel: {
						margin: 20,
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
		    yAxis : [{
		    	type : 'value',
		    	show: false
		    }],
		    series : $scope.StreamPeopledata
		};
		var option3 = {
      tooltip: {
        trigger: 'item',
        formatter: "{b}: {d}%"
      },
      legend: {
        orient: 'vertical',
        left:10,
        padding:[0,0,56,0],
        x: 'left',
        y: 'bottom',
        icon: 'circle',
        itemGap:20,
        itemWidth:40,
        textStyle: {
          color: '#7f8fa4',
          fontSize: 14,
          fontWeight: ''
        },
        data: $scope.inDeviceItems[0].legend,
      },
      series: [
        {
          name:'',
          type:'pie',
          radius: ['14%', '40%'],
          center:['50%','30%'],
          avoidLabelOverlap: false,
          hoverAnimation:false,
          label: {
            normal: {
              show: false,
              position: 'center'
            },
            emphasis: {
              show: false,
              textStyle: {
                fontSize: '30',
                fontWeight: 'bold'
              }
            }
          },
          labelLine: {
            normal: {
              show: false
            }
          },
          data:$scope.inDeviceItems[0].seriesData
        }
      ]
    };
    var option4 = {
        tooltip: {
          trigger: 'item',
          formatter: "{b}: {d}%"
        },
        legend: {
          orient: 'vertical',
          left:10,
          padding:[0,0,20,0],
          x: 'left',
          y: 'bottom',
          icon: 'circle',
          itemGap:20,
          itemWidth:40,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          },
          data: $scope.inDeviceItems[1].legend//['运达小区','测试校区555','ceshi321']
        },
        series: [
          {
            name:'',
            type:'pie',
            radius: ['15%', '40%'],
            center:['50%','30%'],
            avoidLabelOverlap: false,
            hoverAnimation:false,
            label: {
              normal: {
                show: false,
                position: 'center'
              },
              emphasis: {
                show: false,
                textStyle: {
                  fontSize: '30',
                  fontWeight: 'bold'
                }
              }
            },
            labelLine: {
              normal: {
                show: false
              }
            },
            data:$scope.inDeviceItems[1].seriesData
          }
        ]
      };
    var option5 = {
        tooltip: {
          trigger: 'item',
          formatter: "{b}: {d}%"
        },
        legend: {
          orient: 'vertical',
          left:10,
          padding:[0,0,20,0],
          x: 'left',
          y: 'bottom',
          icon: 'circle',
          itemGap:20,
          itemWidth:40,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          },
          data: $scope.inDeviceItems[2].legend
        },
        series: [
          {
            name:'',
            type:'pie',
            radius: ['15%', '40%'],
            center:['50%','30%'],
            avoidLabelOverlap: false,
            hoverAnimation:false,
            label: {
              normal: {
                show: false,
                position: 'center'
              },
              emphasis: {
                show: false,
                textStyle: {
                  fontSize: '30',
                  fontWeight: 'bold'
                }
              }
            },
            labelLine: {
              normal: {
                show: false
              }
            },
            data:$scope.inDeviceItems[2].seriesData
          }
        ]
      };
    var option6 = {
        tooltip: {
          trigger: 'item',
          formatter: "{b}: {d}%"
        },
        legend: {
          orient: 'vertical',
          left:10,
          padding:[0,0,90,0],
          x: 'left',
          y: 'bottom',
          icon: 'circle',
          itemGap:20,
          itemWidth:40,
          textStyle: {
            color: '#7f8fa4',
            fontSize: 14,
            fontWeight: ''
          },
          data: $scope.inDeviceItems[3].legend,
        },
        series: [
          {
            name:'',
            type:'pie',
            radius: ['15%', '40%'],
            center:['50%','30%'],
            avoidLabelOverlap: false,
            hoverAnimation:false,
            label: {
              normal: {
                show: false,
                position: 'center'
              },
              emphasis: {
                show: false,
                textStyle: {
                  fontSize: '30',
                  fontWeight: 'bold'
                }
              }
            },
            labelLine: {
              normal: {
                show: false
              }
            },
            data:$scope.inDeviceItems[3].seriesData
          }
        ]
      };
        // 使用刚指定的配置项和数据显示图表。
    myChart2.setOption(option2);
		myChart3.setOption(option3);
		myChart4.setOption(option4);
		myChart5.setOption(option5);
		myChart6.setOption(option6);
    }

//  $scope.changeSelect = function(){
//  	console.log( $scope.school.schoolItem)
//  }
	Init_load();


	/*   ----------------- wu 2016-12-19 --------------------    */

	$scope.jsdl_dlfx_rlt_option = null;
	//    插件宽高比例问题，，， 4：3， 待验证。


	var geoCoordMap = {

	    "长沙1":[-10,16 * 0.75],
		"长沙2":[-11,13 * 0.75],
		"长沙3":[-12,15 * 0.75],
		"长沙4":[-14,16 * 0.75]

	};

	var convertData = function (data) {
	    var res = [];
	    for (var i = 0; i < data.length; i++) {
	        var geoCoord = geoCoordMap[data[i].name];
	        if (geoCoord) {
	            res.push(geoCoord.concat(data[i].value));
	        }
	    }
	    return res;
	};

	$scope.jsdl_dlfx_rlt_option = {
	    visualMap: {
	        min: 0,
	        max: 1,
	        splitNumber: 0,
	        inRange: {
	            color: ['#d94e5d','#eac736','#50a3ba'].reverse()
	        },
	        textStyle: {
	            color: '#fff'
	        },
	        show:false
	    },
	    geo: {
	        map: 'xinyun18',
	        label: {
	            emphasis: {
	                show: false
	            }
	        },
	        roam: false,
	        itemStyle: {
	            normal: {
	                areaColor: 'transparent',
	                borderColor: 'transparent'
	            },
	            emphasis: {
	                areaColor: 'transparent'
	            }
	        },
	        width:'100%'


	    },
	    series: [{
	        name: 'AQI',
	        type: 'heatmap',
	        coordinateSystem: 'geo',
	        data: convertData([

	            {name: "长沙1", value: 0.77},
	            {name: "长沙2", value: 0.91},
				{name: "长沙3", value: 0.7},
				{name: "长沙4", value: 1}
	        ])
	    }]
	};


	//    人员定位  ------------------------------------------------
	//   option
	$scope.jsdl_dlfx_dw_option = null;
	//   定位数据
	$scope.dw_data = [];


	//    插件宽高比例问题，，， 4：3， 待验证。

	$scope.jsdl_dlfx_dw_option = {
		//   用于触发修改监听
		chartOptionBul:true,
	    tooltip : {
	        trigger: 'item'
	    },
	    geo: {
	        map: 'xinyun18',
	        label: {
	            emphasis: {
	                show: false
	            }
	        },
	        roam: false,
	        itemStyle: {
	            normal: {
	                areaColor: 'transparent',
	                borderColor: 'transparent'
	            },
	            emphasis: {
	                areaColor: 'transparent'
	            }
	        },
	        width:'100%'
	    },
	    series : [
	        {
	            name: '定位',
	            type: 'scatter',
	            coordinateSystem: 'geo',
	            data:$scope.dw_data,
	            symbolSize: function (val) {
	                return val[2] / 10;
	            },
	            label: {
	                normal: {
	                    formatter: '{b}',
	                    position: 'right',
	                    show: false
	                },
	                emphasis: {
	                    show: true
	                }
	            },
	            itemStyle: {
	                normal: {
	                    color: '#ddb926'
	                }
	            }
	        },
	        {
	            name: '定位',
	            type: 'effectScatter',
	            coordinateSystem: 'geo',
	            data:$scope.dw_data,
	            symbolSize: function (val) {
	                return val[2] / 10;
	            },
	            showEffectOn: 'render',
	            rippleEffect: {
	                brushType: 'stroke'
	            },
	            hoverAnimation: true,
	            label: {
	                normal: {
	                    formatter: function(a){
	                    	return a.name;
	                    },
	                    position: 'right',
	                    show: true
	                }
	            },
	            itemStyle: {
	                normal: {
	                    color: '#f4e925',
	                    shadowBlur: 10,
	                    shadowColor: '#333'
	                }
	            },
	            zlevel: 1
	        }
	    ]
	};



	//   取实时人员坐标
    var getclassroompeopleinfo = function(){
    	var url=config.HttpUrl+"/basicset/getclassroompeopleinfo";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Begindate":"2016-11-28 00:00:00",
            "Enddate":"2016-12-28 23:59:59",
			"Classroomid": 343
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取实时人员坐标",data)
			if(data.Rcode == "1000") {
				$scope.geoCoordMap = {};
				$scope.dw_data = [];
				for(var a in data.Result){
					if(data.Result[a].X != 0 && data.Result[a].Y != 0)$scope.dw_data.push({"value":[data.Result[a].X.toFixed(3),(data.Result[a].Y * 0.75).toFixed(3),150],"name":data.Result[a].Truename});
				}
				//
				$scope.jsdl_dlfx_dw_option.series[0].data = $scope.dw_data;
				$scope.jsdl_dlfx_dw_option.series[1].data = $scope.dw_data;
				//   触发 修改监听
				$scope.jsdl_dlfx_dw_option.chartOptionBul = !($scope.jsdl_dlfx_dw_option.chartOptionBul);
            }else{
              toaster.pop('warning',data.Reason);
            }
		}, function(reason) {}, function(update) {});
    }

	getclassroompeopleinfo();

	/*   ----------------- wu 2016-12-19 End --------------------    */


}]);
