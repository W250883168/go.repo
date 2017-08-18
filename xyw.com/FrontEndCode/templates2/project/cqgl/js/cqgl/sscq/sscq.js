'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 实时出勤
 */

/*   出勤统计-实时出勤      */
app.controller('cqglSscqContr', ['$scope', 'httpService','$interval','$state','$modal','toaster', function($scope, httpService,$interval,$state,$modal,toaster) {
	console.log("实时出勤")
	//   默认校区  教学楼  本地存储
	$scope.defaultSchoolFloor = {
		"school": [],
		"floor": []
	};


	//开始定义定时器
	var tm=$scope.setglobaldata.gettimer("sscq");
	if(tm.Key!="sscq"){
		tm.Key="sscq";
		tm.keyctrl="app.cqgl.sscq";
		tm.fnAutoRefresh=function(){
			console.log("开始调用定时器");
			this.interval = $interval(function() {
				buildingString();//   格式化字符串
				getClassroomStatusList(buildingids);//   get 教室
			}, config.sscqRefreshTime);
		};
		tm.fnStopAutoRefresh=function(){
			console.log("进入取消方法");
			if(!angular.isUndefined(this.interval)) {
				$interval.cancel(this.interval);
				this.interval = 'undefined';
				console.log("进入取消成功");
			}
			this.interval=null;
		};
		$scope.setglobaldata.addtimer(tm);
	}
	//结束定义定时器


	//   取校区
	$scope.schoolItems = [];
	var getcampus = function() {
		var url = config.HttpUrl + "/basicset/getcampus";
		var data = {

		};
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {

				$scope.schoolItems = data.Result;
				//   加入默认选中项
				if($scope.defaultSchoolFloor["school"].length > 0) {
					for(var i = 0; i < $scope.schoolItems.length; i++) {
						if($scope.schoolItems[i].Campusid == $scope.defaultSchoolFloor["school"][0].Campusid) {
							$scope.schoolItems[i].checkbox = true;
						}
					}
					getcampusFloor($scope.defaultSchoolFloor["school"][0].Campusid);
				} else {
					$scope.schoolItems[0].checkbox = true;
					getcampusFloor($scope.schoolItems[0].Campusid);
				}
				//    定时器
				tm.fnAutoRefreshfn(tm);
			} else {
        toaster.pop('warning',data.Reason);
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}
	getcampus();

	//  取楼栋
	$scope.floorItems = [];
	var getcampusFloor = function(campusid) {


		var url = config.HttpUrl + "/basicset/getbuilding";
		var data = {
			//"Usersid": config.GetUser().Usersid,
			//"Rolestype": config.GetUser().Rolestype,
			//"Token": config.GetUser().Token,
			//"Os": "WEB",
			"campusid": campusid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.floorItems = data.Result;
				//   加入默认选中项
				for(var i = 0; i < $scope.floorItems.length; i++) {
					if($scope.defaultSchoolFloor["floor"].length>0){
					for(var b = 0; b < $scope.defaultSchoolFloor["floor"].length; b++) {
						if($scope.floorItems[i].Buildingid == $scope.defaultSchoolFloor["floor"][b].Buildingid) {
							$scope.floorItems[i].checkbox = true;
						}
					}
					}else{
						$scope.floorItems[i].checkbox = true;
					}
				}
				if($scope.defaultSchoolFloor["floor"].length==0){
					$scope.floor_checkbox();
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}

	//   校区选择
	$scope.school_tab = function(item) {
		for(var i = 0; i < $scope.schoolItems.length; i++) {
			if($scope.schoolItems[i].Campusid == item.Campusid) {
				$scope.schoolItems[i].checkbox = true;
			} else {
				$scope.schoolItems[i].checkbox = false;
			}
		}
		//清除楼栋选择数据
		$scope.defaultSchoolFloor["floor"]=[];
		if(localStorage.getItem("jsjkTab"+item.Campusid)!=undefined){
		var defaultSchoolFloorlib = JSON.parse(localStorage.getItem("jsjkTab"+item.Campusid));
		$scope.defaultSchoolFloor=defaultSchoolFloorlib;
		}
		//   查楼栋
		getcampusFloor(item.Campusid);

		if($scope.defaultSchoolFloor.school[0] != undefined){
			if(item.Campusid == $scope.defaultSchoolFloor.school[0].Campusid){
				//   格式化字符串
				buildingString();
				//   get 教室
				getClassroomStatusList(buildingids);
			}else{
				$scope.classroomItems = [];
			}
		}
	}

	//   楼栋选择
	$scope.floor_checkbox = function(item) {

		//   校区
		$scope.defaultSchoolFloor["school"] = [];
		for(var i = 0; i < $scope.schoolItems.length; i++) {
			if($scope.schoolItems[i].checkbox) {
				$scope.defaultSchoolFloor["school"].push($scope.schoolItems[i]);
			}
		}
		//   楼栋
		$scope.defaultSchoolFloor["floor"] = [];
		for(var i = 0; i < $scope.floorItems.length; i++) {
			if($scope.floorItems[i].checkbox) {
				$scope.defaultSchoolFloor["floor"].push($scope.floorItems[i]);
			}
		}

		//  格式化
		var temp = null;
		temp = $scope.defaultSchoolFloor["school"][0];
		$scope.defaultSchoolFloor["school"][0] = {};
		$scope.defaultSchoolFloor["school"][0].Campusid = temp.Campusid;
		$scope.defaultSchoolFloor["school"][0].Campusname = temp.Campusname;
		for(var i = 0; i < $scope.defaultSchoolFloor["floor"].length; i++) {
			temp = $scope.defaultSchoolFloor["floor"][i];
			$scope.defaultSchoolFloor["floor"][i] = {};
			$scope.defaultSchoolFloor["floor"][i].Buildingid = temp.Buildingid;
			$scope.defaultSchoolFloor["floor"][i].Buildingname = temp.Buildingname;
		}
		localStorage.setItem("jsjkTab"+$scope.defaultSchoolFloor["school"][0].Campusid, JSON.stringify($scope.defaultSchoolFloor));
		localStorage.setItem("jsjkTab", JSON.stringify($scope.defaultSchoolFloor));
//console.log(localStorage.getItem("jsjkTab"));

		//   格式化字符串
		buildingString();
		//   get 教室
		getClassroomStatusList(buildingids);
	}



	/////////////////////////////////////////////////////////

	//------------------------------------------------------------------------------------------------------------------
    var fnTransformData = function (data){
    	var d;
        d = data;

        //定义变量
        var bid = undefined;    //楼栋id
        var fid = undefined;    //楼层id
        var ob = undefined;     //object of building
        var of = undefined;     //object of floor
        var oc = undefined;     //object of classroom
        var xh = 0;             //序号
        var maxCol = 0;         //所有楼层中的最大教室数量（所有楼层按这个数量显示教室列数）
        var rd=[];              //存放最后的数据

        //循环对数据进行处理,构造更易在界面上展现的数据集
        for (var i=0;i< d.length;i++){
        	var ir2 = null;
            xh++; //序号直接加1(遇新楼层时，恢复为1

            //取出教室编号后两位
            var name = d[i].ClassroomName;
            ir2 = parseInt(name.substring(name.length-2,name.length));// 将编号右边两位转换为整数

            //oc
            oc = {ClassroomId:d[i].ClassroomId,
                ClassroomName:d[i].ClassroomName,
                ClassroomState:d[i].ClassroomState,
                Sumnumbers:d[i].Sumnumbers,
                Seatsnumbers:d[i].Seatsnumbers,
                Classroomstype:d[i].Classroomstype,
                Percentage:1
            };
            var Percentages= (Number(oc.Sumnumbers)/Number(oc.Seatsnumbers));
            Percentages=Number(Percentages);
            //判断教室状态
            if (oc.ClassroomState==-1) {//教室离线
                oc.Percentage = -1;
            }
            else if (oc.ClassroomState==0) {//教室开放
                /*if (Percentages >= 0.5 && Percentages <= 0.8) {//教室内人数正常
                    oc.Percentage = 2;
                }
                else if (Percentages > 0.8) {//教室内人数比较多
                    oc.Percentage = 3;
                }
                else if(Percentages < 0.5 && Percentages > 0) {//教室内人少
                    oc.Percentage = 1;
                }else{
                	oc.Percentage = 0;
                }*/
            	oc.Percentage = 0;
            }
            else if (oc.ClassroomState==1) {//上课中
                if (Percentages >= 0.5 && Percentages <= 0.8) {//教室内人数正常
                    oc.Percentage = 2;
                }
                else if (Percentages > 0.8) {//教室内人数比较多
                    oc.Percentage = 3;
                }
                else{
                	oc.Percentage = 1;
                }
            }
            //of
            if (d[i].FloorId != fid){
                //将当前楼层id保存到fid
                fid = d[i].FloorId;

                //将前一个楼层的最大教室数保存起来
                if (xh-1>maxCol){
                    maxCol = xh-1
                }

                //保存前一个
                if (of != undefined){
                    ob.data.push(of);
                    of = undefined;
                }

                //建立新的楼层
                of = {FloorId:d[i].FloorId,FloorName:d[i].FloorName,FloorImage:d[i].FloorsImage,data:[]};
                xh = 1;//建立新的楼层后，教室序号肯定是从1开始，所以初始化为1
            }

            //将教室压入楼层前,先补插空缺教室（按顺序补齐后台数据中没有返回的教室）
            for (var k=0;k<ir2-xh;k++){
                of.data.push(fnGetNullClassroom(name.substring(0,name.length-2),xh)) //
                xh++;//将序号+1
            }

            //将教室压入楼层
            of.data.push(oc)

            //ob
            if (d[i].BuildingId != bid){
                //保存前一个
                if (ob != undefined){
                    rd.push(ob);    //将ob压入rd
                    ob = undefined;
                }
                bid = d[i].BuildingId;

                //建立新的楼栋
                ob = {BuildingId:d[i].BuildingId,BuildingName:d[i].BuildingName,data:[]};
            }

        }

        //最后一栋
        if (ob != undefined){
            if (of != undefined){
                ob.data.push(of)//最后一层压入楼栋
            }
            rd.push(ob);//将ob压入rd
        }

        //将前一个楼层的最大教室数保存起来
        if (xh>maxCol){
            maxCol = xh
        }

        //返回
        return {MaxCol:maxCol,data:rd}
    }

    //获得空教室对象
    var fnGetNullClassroom = function(floorCode,classroomName){
        //将教室编码补齐
        var fullClassroomName
        if (classroomName<10){
            fullClassroomName = floorCode + "0" + classroomName
        }else{
            fullClassroomName = floorCode + classroomName
        }
        //凡是补插的空教室，教室id一律为-1
        return {ClassroomId:-1,ClassroomName:fullClassroomName,ClassroomState:-1,CollectionNumbers:0,HaveStop:-1,HaveAlert:-1,HaveOffline:-1,HaveRun:-1}
    }

    //增加楼层教室（补齐楼层教室）
    var fnAddFloorClassroom = function (data){
        //以所有楼层最大教室数量为参考，将楼层教室数量小于最大数量的，补齐
        var rd = data.data;
        var maxCol = data.MaxCol;
        for (var i=0;i<rd.length;i++){
            var f = rd[i].data;
            for (var j=0;j<f.length;j++){
                var c = f[j].data;
                //如果当前楼层的教室数量小于maxCol，则循环补齐
                if (c.length<maxCol){
                    //得到教室名称的前缀(使用第一个教室的名称，除去后两位即可)
                    name = c[0].ClassroomName;
                    //当前楼层已有教室数量
                    var len = c.length;
                    //循环补齐
                    for (var k=0;k<maxCol-len;k++){
                        c.push(fnGetNullClassroom(name.substring(0,name.length-2),len+1+k))
                    }
                }
            }
        }

        return {MaxCol:maxCol,data:rd};
    }
    //------------------------------------------------------------------------------------------------------------------
	/////////////////////
	//   取教室列表
	//  被选中楼栋ID
	var buildingids = "";
	var buildingString = function(){
		//   楼栋
		buildingids = "";
		for(var i = 0; i < $scope.defaultSchoolFloor.floor.length; i++) {
			//   楼栋ID
			buildingids += $scope.defaultSchoolFloor.floor[i].Buildingid + ",";
		}
		if(buildingids.length > 0) {
			buildingids = buildingids.substr(0, buildingids.length - 1);
		}else{
			buildingids = "";
		}
	}

	//   设备控制 取教室列表
	$scope.classroomItems = [];
	var getClassroomStatusList = function(buildingids) {
		var url = config.HttpUrl + "/action/getattendancelist";
		var data = {
            //用户id：整型
            "Usersid": config.GetUser().Usersid,
            //角色类型：整型
            "Rolestype": config.GetUser().Rolestype,
            //令牌：字符串
            "Token": config.GetUser().Token,
            //操作系统：字符串
            "Os": "WEB",
            //多个楼栋id，逗号拼接成串
			Buildingids: buildingids,
			"State":1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.classroomItems = fnTransformData(data.Result).data;
				console.log($scope.classroomItems)
				//$scope.classroomItems = fnAddFloorClassroom($scope.classroomItems).data;
				//console.log($scope.classroomItems)
			} else {
				$scope.classroomItems = [];
        toaster.pop('warning', data.Reason);
			}
			console.log('设备控制 取教室列表',$scope.classroomItems)
		}, function(reason) {}, function(update) {});
	}

	$scope.openSref = function(item){
		//    上课中
		if(item.ClassroomState == 1){
			$state.go("app.cqgl.sscq.details",{"ClassroomId":item.ClassroomId});
		}
	}


	//   载入默认教室
	if(localStorage.getItem("jsjkTab") != null) {
		$scope.defaultSchoolFloor = JSON.parse(localStorage.getItem("jsjkTab"));
		//   格式化字符串
		buildingString();
		//   get 教室
		getClassroomStatusList(buildingids);
	}else{

		}

	//   查看楼层平面图
	$scope.openPic = function(floor) {
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbjk/modal_pic.html',
			controller: 'picSbjkContr',
			windowClass: 'm-sbjk-modal',
			//size: "lg",
			resolve: {
				items: function() {
					return floor;
				}
			}
		});
	}

}]);




/*   实时出勤查看     */
app.controller("cqglSscqDetailsContr",['$scope','$location','httpService','$filter','toaster',function($scope,$location,httpService,$filter,toaster){
	console.log("实时出勤查看");
	//
	$scope.form = {
		//  开始时间
		"date_begin":"",
		//  结束时间
		"date_end":"",
		//   教室id
		"classroomid": $location.search().ClassroomId,
		//------------------//
		//   上课日期
		"day":"",
		//    上课时间
		"time":"",
		//    教室
		"room":"",
		//   上课老师
		"teacher":"",
		//    课程
		"lesson":"",
		//---------------------//
		//    上课班级
		"class_room":"",
		//    应到人数
		"due_people":"",
		//    实到人数
		"actual_people":"",
		//    缺勤人数
		"lack_people":"",
		//    出勤率
		"cq_people":""
	}

	//  教室信息
	$scope.classRoomItem = "";
	//
	$scope.pointtos = [];
	//
	$scope.schedule = "";

	//   开始时间
	$scope.getBeginDate = function(){
		jeDate({
			dateCell: "#sscq_d_begin",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				console.log(val)
				$scope.$apply(function(){
					$scope.form.date_begin = val;
				});
			},
			okfun: function(elem,val) {
				console.log(val)
				$scope.$apply(function(){
					$scope.form.date_begin = val;
				});
			},
			clearfun:function(elem, val) {
				console.log(val)
				$scope.$apply(function(){
					$scope.form.date_begin = "";
				});
			}
		});
	}
	//   结束时间
	$scope.getEndDate = function(){
		jeDate({
			dateCell: "#sscq_d_end",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.$apply(function(){
					$scope.form.date_end = val;
				});
			},
			okfun: function(elem,val) {
				$scope.$apply(function(){
					$scope.form.date_end = val;
				});
			},
			clearfun:function(elem, val) {
				$scope.$apply(function(){
					$scope.form.date_end = "";
				});
			}
		});
	}


	//    取教室信息
	var getClassroomInfo = function(classroomid) {
		if(Number(classroomid) < 0 && !classroomid){return false}else{classroomid = Number(classroomid)};
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			id: classroomid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			console.log('取教室信息',data)
			if(data.Rcode == "1000") {
				$scope.classRoomItem = data.Result;
				//   上课班级

				$scope.form.class_room = $scope.classRoomItem.Classesname;
				//   上课老师 Truename
				$scope.form.teacher = $scope.classRoomItem.Truename;
				//   课程 Chaptername
				$scope.form.lesson = $scope.classRoomItem.Curriculumname;
				//   教室
				$scope.form.room = $scope.classRoomItem.Buildingname + "-" + $scope.classRoomItem.Campusname + "-" + $scope.classRoomItem.Floorname + "-" + $scope.classRoomItem.Classroomsname;
				var Curriculumclassroomchaptercentreids="";
				if($scope.classRoomItem.Qccci!=null){
					$scope.form.class_room="";
					if($scope.classRoomItem.Qccci.length>0){
						for(var i=0;i<$scope.classRoomItem.Qccci.length;i++){
							$scope.form.class_room=$scope.form.class_room+"   "+$scope.classRoomItem.Qccci[i].Classesname;
							Curriculumclassroomchaptercentreids=Curriculumclassroomchaptercentreids+$scope.classRoomItem.Qccci[i].Curriculumclassroomchaptercentreid+",";
						}
					}
					Curriculumclassroomchaptercentreids=Curriculumclassroomchaptercentreids.substring(0,Curriculumclassroomchaptercentreids.length-1);
					console.log((Curriculumclassroomchaptercentreids!=""));
					console.log((Curriculumclassroomchaptercentreids!="undefined"));
					if(Curriculumclassroomchaptercentreids!="" && Curriculumclassroomchaptercentreids!="undefined"){
						getPointtos(Curriculumclassroomchaptercentreids);
					}else{
						getPointtos($scope.classRoomItem.Curriculumclassroomchaptercentreid);
					}
				}
				getCurriculums($scope.form.classroomid);
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    get 点到学生信息
	var getPointtos = function(id) {
		console.log(id);
		console.log("--------------------------");
		if(id!="" && id!="undefined" && !id){return false}else{};
		//   不在上课时间 id 为0
//		if(id == 0){return false}
		var url = config.HttpUrl + "/action/getpointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
//			"Curriculumclassroomchaptercentreid": id
			"Curriculumclassroomchaptercentreids":id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.pointtos = data.Result;
				//   实到 与 未到
				if($scope.pointtos != null){
					if($scope.pointtos.length > 0){
						var temp = 0;
						for(var a in $scope.pointtos){
							if($scope.pointtos[a].State == 0){
								temp += 1;
							}
						}
						//   应到人数
						$scope.form.due_people = $scope.pointtos.length;
						//   实到
						$scope.form.actual_people = $scope.form.due_people - temp;
						//   未到
						$scope.form.lack_people = temp;
						//   出勤率
						$scope.form.cq_people = Math.round($scope.form.actual_people / $scope.form.due_people * 1000) / 10.0 + "%";
					}
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//  用户查课表  当天内
	var getCurriculums = function(classroomid) {
		if(Number(classroomid) < 0 && !classroomid){return false}else{classroomid = Number(classroomid)};

		var myDate = new Date();
		var day =  $filter('date')(myDate, 'yyyy-MM-dd');

		var url = config.HttpUrl + "/action/getcurriculums";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": day + " 00:00:01",
			"Enddate": day + " 23:59:59",
			"State": -1,
			//"Teacherids":config.GetUser().Usersid.toString(),
			"Teacherids":"",
			"Classroomid":classroomid,
			//"PageSize":70,
			"PageIndex":-1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.schedule = data.Result;
				//   上课日期
				if($scope.schedule.length > 0){
					var temp = {};
					for(var a in $scope.schedule){
						if($scope.schedule[a].Curriculumclassroomchaptercentreid == $scope.classRoomItem.Curriculumclassroomchaptercentreid){
							temp = $scope.schedule[a];
							break;
						}
					}
					//   上课日期
					if(temp != null && !angular.equals({},temp)){
						$scope.form.day = temp.Begindate.substr(0,10);
						$scope.form.time = temp.Begindate.substr(11,5) + "-" + temp.Enddate.substr(11,5);
					}
				}


			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}


	//   run
	var run = function(){
		//   当天
		var myDate = new Date();
		var day =  $filter('date')(myDate, 'yyyy-MM-dd');

		$scope.form.date_begin = day + " 00:00:00";
		$scope.form.date_end = day + " 23:59:59";

		getClassroomInfo($scope.form.classroomid);
	}
	run();
}]);


/*    实时出勤-教室-查找        */
app.controller("cqglSscqDetailsListContr",['$scope','$location','httpService','$filter','toaster',function($scope,$location,httpService,$filter,toaster){
	console.log('实时出勤-教室-查找')
	$location.search().ClassroomId;

	$scope.form = {
		"ClassroomId":$location.search().ClassroomId,
		"begin_date":$location.search().be,
		"end_date":$location.search().end
	}





	console.log($scope.form)



	/*  -------------------- 分页、页码  -----------------------  */
	$scope.backPage = {
		"PageCount":"",
		"PageIndex":"1",
		"PageSize":"15",
		"RecordCount":""
	};
	/*----------------
	//    分页对象添加页码
	//    return  obj  分页对象
	//    pagedata:obj  分页对象
	//    maxpagenumber:int  显示页码数默认5个页码
	------------------*/
	var pageFn = function(pagedata,maxpagenumber){
		if(pagedata.length < 1)return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber){
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1 ){
				Snumber = 1
			}else if((pagedata.PageIndex + mid) > pagedata.PageCount){
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			}else{
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++){
				nub.push(Snumber + i);
			}
		}else{
			for(var i = 0; i < pagedata.PageCount; i++){
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	}


	//  翻页
	$scope.pageClick = function(pageindex){
		if(!(Number(pageindex) > 0))return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount){
			$scope.backPage.PageIndex=pageindex;
			getCurriculums();
		}
	}
	/*  -------------------- 分页、页码  -----------------------  */

	//   查询
	$scope.searchPost = function(){
		$scope.backPage.PageIndex=1;
		getCurriculums();
	}


	//   开始时间
	$scope.getBeginDate = function(){
		jeDate({
			dateCell: "#lscq_begin",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.begin_date = val;
				//
				$scope.searchPost();
			},
			okfun: function(elem,val) {
				$scope.form.begin_date = val;
				//
				$scope.searchPost();
			},
			clearfun:function(elem, val) {
				$scope.form.begin_date = "";
			}
		});
	}
	//   结束时间
	$scope.getEndDate = function(){
		jeDate({
			dateCell: "#lscq_end",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			maxDate: jeDate.now(0),
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.end_date = val;
				//
				$scope.searchPost();
			},
			okfun: function(elem,val) {
				$scope.form.end_date = val;
				//
				$scope.searchPost();
			},
			clearfun:function(elem, val) {
				$scope.form.end_date = "";
			}
		});
	}

	//  用户查课表  当天内
	var getCurriculums = function() {
		var myDate = new Date();
		var temp_end = "";
		new Date($scope.form.end_date) > myDate ? temp_end = $filter('date')(myDate, 'yyyy-MM-dd HH:mm:ss') : temp_end = $scope.form.end_date;

		var url = config.HttpUrl + "/action/gethistoryattendance";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": $scope.form.begin_date,
			"Enddate": temp_end,
			"State": -1,
			"Classroomid":Number($scope.form.ClassroomId),
			"PageSize":Number($scope.backPage.PageSize),
			"PageIndex":Number($scope.backPage.PageIndex)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.schedule = data.Result.PageData;
				//   出勤率
				for(var a in $scope.schedule){
					$scope.schedule[a].Toclassrate = Math.round($scope.schedule[a].Toclassrate * 1000) / 10.0 + "%";
				}

				//   分页
				var objPage={PageCount:0,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				//   分页
				//$scope.backPage.PageIndex = $scope.schedule;
				$scope.backPage = pageFn(objPage,5);
			} else {
				$scope.schedule = [];
        toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}


	//    取教室信息
	var getClassroomInfo = function(classroomid) {
		if(Number(classroomid) < 0 && !classroomid){return false}else{classroomid = Number(classroomid)};
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			id: classroomid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.classRoomItem = data.Result;
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   run
	var run = function(){
		//   当天
		//var myDate = new Date();
		//var day =  $filter('date')(myDate, 'yyyy-MM-dd');

		//$scope.form.begin_date = day + " 00:00:00";
		//$scope.form.end_date = day + " 23:59:59";


		getClassroomInfo($scope.form.ClassroomId);
		getCurriculums();

	}
	run();

}]);
