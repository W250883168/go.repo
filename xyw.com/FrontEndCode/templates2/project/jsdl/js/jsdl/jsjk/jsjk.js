'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

/*   教室监控     */
app.controller('jsdlJsjkContr', ['$scope','$rootScope', 'httpService', '$modal','$interval','$state','toaster', function($scope, $rootScope,httpService, $modal,$interval,$state,toaster) {
	console.log("教室监控")
	//   默认校区  教学楼  本地存储
	$scope.defaultSchoolFloor = {
		"school": [],
		"floor": []
	};
//开始定义定时器
	var tm=$scope.setglobaldata.gettimer("jsjk");
	if(tm.Key!="jsjk"){
	tm.Key="jsjk";
	tm.keyctrl="app.jsdl.jsjk";
	tm.fnAutoRefresh=function(){
		this.interval = $interval(function() {
			buildingString();//   格式化字符串
			getClassroomStatusList(buildingids);//   get 教室
		}, config.jsjkRefreshTime);
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
				//   自动刷新
				var tm=$scope.setglobaldata.gettimer("jsjk");
				tm.fnAutoRefreshfn(tm);
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
                if (Percentages >= 0.3 && Percentages < 0.7) {//教室内人数正常
                    oc.Percentage = 2;
                }
                else if (Percentages >= 0.7) {//教室内人数比较多
                    oc.Percentage = 3;
                }
                else {//教室内人少
                    oc.Percentage = 1;
                }
            }
            else if (oc.ClassroomState==1) {//上课中
                /*if (Percentages >= 0.3 && Percentages < 0.7) {//教室内人数正常
                    oc.Percentage = 2;
                }
                else if (Percentages >= 0.7) {//教室内人数比较多
                    oc.Percentage = 3;
                }
                else {//教室内人少
                    oc.Percentage = 1;
                }*/
                oc.Percentage = 0;
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
		var url = config.HttpUrl + "/basicset/getpeoples";
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
			Buildingids: buildingids
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
        toaster.pop('warning',data.Reason);
			}
			//console.log($scope.classroomItems)
		}, function(reason) {}, function(update) {});
	}

	//    打开
	$scope.openSref = function(item){
		//    开放中
		if(item.ClassroomState == 0){
			$state.go("app.jsdl.jsjk.sstj",{"ClassroomId":item.ClassroomId});
		}
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



	//   载入默认教室
	if(localStorage.getItem("jsjkTab") != null) {
		$scope.defaultSchoolFloor = JSON.parse(localStorage.getItem("jsjkTab"));
		//   格式化字符串
		buildingString();
		//   get 教室
		getClassroomStatusList(buildingids);
	}else{

		}
//   教室导流 教室详情 教室内人员列表信息查询
    var LoadClassroomPeopleInfo = function() {

        var url = config.HttpUrl + "/curriculum/getcurriculumchaptersinfo";
        var data = {
			"Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Teacherid": 26,
            "Classesid":1
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function(data) {
            if(data.Rcode == "1000") {
			console.log(data.Result);
            } else {
            }
        }, function(reason) {}, function(update) {});
    };
    LoadClassroomPeopleInfo();
	$scope.$on('to-parent', function(event,data) {
        console.log('ParentCtrl', data);       //父级能得到值
   });
}]);
