'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理-查课表
 */

app.filter('propsFilter', function() {
	return function(items, props) {
		var out = [];
		if(angular.isArray(items)) {
			items.forEach(function(item) {
				var itemMatches = false;
				var keys = Object.keys(props);
				for(var i = 0; i < keys.length; i++) {
					var prop = keys[i];
					try {
						if(typeof(props[prop]) == "number") {
							var text = props[prop];
							if(item[prop].toString() == text.toString()) {
								itemMatches = true;
								break;
							}
						} else {
							var text = props[prop].toString().toLowerCase();
							if(item[prop].toString().toLowerCase().indexOf(text) !== -1) {
								itemMatches = true;
								break;
							}
						}
					} catch(e) {}
				}
				if(itemMatches) {
					out.push(item);
				}
			});
		} else {
			out = items;
		}
		return out;
	};
});

/*    课程管理-查课表      */
app.controller("kcglCkbContr", ['$scope', '$state', 'httpService', '$filter','toaster', function($scope, $state, httpService, $filter,toaster) {
	//$state.go("app.qxgl",false);
	console.log("课程管理-查课表")

	$scope.Campus = {};
	$scope.Campus.Campuslist = [];
	$scope.Campus.selectCampusitem = "";

	$scope.building = {};
	$scope.building.buildinglist = [];
	$scope.building.selectbuildingitem = "";

	$scope.floors = {};
	$scope.floors.floorslist = [];
	$scope.floors.selectfloorsitem = "";

	$scope.college = {};
	$scope.college.collegelist = [];
	$scope.college.selectcollegeitems = [];

	$scope.major = {};
	$scope.major.majorlist = [];
	$scope.major.selectmajoritems = [];

	$scope.classe = {};
	$scope.classe.classelist = [];
	$scope.classe.selectclasseitems = [];

	$scope.Teacher = {};
	$scope.Teacher.Teacherlist = [];
	$scope.Teacher.selectteacheritem = [];
	//   分页
	$scope.backPage = {
			"PageIndex": 1,
			"PageSize": 10
		}
		//   表单
	$scope.form = {
			"Begindatestr": "",
			"Endatestr": "",
			"Search": ""
		}
		//    查询课表列表
	$scope.getcurriculumsList = [];
	//
	$scope.getAll = [];

	$scope.changeselect = function(index) {
		console.log('1111111111111111111',$scope.Campus.selectCampusitem);
		switch(index) {
			case 0:
				$scope.building.selectbuildingitem = {}; //清除选中的楼栋
			case 1:
				$scope.floors.selectfloorsitem = {};
				break;
			default:
				$scope.Campus.selectCampusitem = {};
				$scope.building.selectbuildingitem = {};
				$scope.floors.selectfloorsitem = {};
		}
	};


	//   学院变化
	$scope.$watch('college.selectcollegeitems',function(newValue,oldValue, scope){
		if(newValue){
			//    专业
			$scope.major.majorlist = [];
			for(var a in newValue){
				for(var b in $scope.getAll[4]){
					if(newValue[a].Id == $scope.getAll[4][b].Collegeid){
						$scope.major.majorlist.push($scope.getAll[4][b]);
					}
				}
			}
			/*//    班级
			$scope.classe.classelist = [];
			for(var a in newValue){
				for(var b in $scope.getAll[5]){
					if(newValue[a].Id == $scope.getAll[5][b].Majorid){
						$scope.classe.classelist.push($scope.getAll[5][b]);
					}
				}
			}*/
			//    老师
			$scope.Teacher.Teacherlist = [];
			for(var a in newValue){
				for(var b in $scope.getAll[6]){
					if(newValue[a].Id == $scope.getAll[6][b].Collegeid){
						$scope.Teacher.Teacherlist.push($scope.getAll[6][b]);
					}
				}
			}
		}

	});

	//   专业变化
	$scope.$watch('major.selectmajoritems',function(newValue,oldValue, scope){
		if(newValue){
			//    班级
			$scope.classe.classelist = [];
			for(var a in newValue){
				for(var b in $scope.getAll[5]){
					if(newValue[a].Id == $scope.getAll[5][b].Majorid){
						$scope.classe.classelist.push($scope.getAll[5][b]);
					}
				}
			}
			//    老师
			$scope.Teacher.Teacherlist = [];
			for(var a in newValue){
				for(var b in $scope.getAll[6]){
					if(newValue[a].Id == $scope.getAll[6][b].Majorid){
						$scope.Teacher.Teacherlist.push($scope.getAll[6][b]);
					}
				}
			}
		}

	});




	var init_data = function() {
		var url = config.HttpUrl + "/basicset/getall";
		var data = {};
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			console.log("getapp",data);
			if(data.Rcode == "1000") {
				$scope.getAll = data.Result;

				$scope.Campus.Campuslist = data.Result[0];
				$scope.building.buildinglist = data.Result[1];
				$scope.floors.floorslist = data.Result[2];
				$scope.college.collegelist = data.Result[3];
				$scope.major.majorlist = data.Result[4];
				$scope.classe.classelist = data.Result[5];
				$scope.Teacher.Teacherlist = data.Result[6];
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//    开始时间
	$scope.showFromDate = function() {
		jeDate({
			dateCell: "#begindate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.Begindatestr = val;
			},
			okfun: function(elem, val) {
				$scope.form.Begindatestr = val;
			},
			clearfun: function(elem, val) {
				$scope.form.Begindatestr = "";
			}
		});
	}

	//    结束时间
	$scope.showToDate = function() {
		jeDate({
			dateCell: "#enddate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.Endatestr = val;
			},
			okfun: function(elem, val) {
				$scope.form.Endatestr = val;
			},
			clearfun: function(elem, val) {
				$scope.form.Endatestr = "";
			}
		});
	}

	$scope.GetselectTeacherItem = function() {
		var selectitemstr = "";
		var k = 0;
		var count = $scope.Teacher.selectteacheritem.length;
		for(k = 0; k < count; k++) {
			selectitemstr = selectitemstr + $scope.Teacher.selectteacheritem[k].Usersid + ","
		}
		selectitemstr = selectitemstr.substring(0, selectitemstr.length - 1);
		return selectitemstr;
	};
	$scope.GetselectcollegeItems = function() {
		var selectitemstr = "";
		var k = 0;
		var count = $scope.college.selectcollegeitems.length;
		for(k = 0; k < count; k++) {
			selectitemstr = selectitemstr + $scope.college.selectcollegeitems[k].Id + ","
		}
		selectitemstr = selectitemstr.substring(0, selectitemstr.length - 1);
		return selectitemstr;
	};
	$scope.GetselectmajorItems = function() {
		var selectitemstr = "";
		var k = 0;
		var count = $scope.major.selectmajoritems.length;
		for(k = 0; k < count; k++) {
			selectitemstr = selectitemstr + $scope.major.selectmajoritems[k].Id + ","
		}
		selectitemstr = selectitemstr.substring(0, selectitemstr.length - 1);
		return selectitemstr;
	};
	$scope.GetselectclasseItems = function() {
		var selectitemstr = "";
		var k = 0;
		var count = $scope.classe.selectclasseitems.length;
		for(k = 0; k < count; k++) {
			selectitemstr = selectitemstr + $scope.classe.selectclasseitems[k].Id + ","
		}
		selectitemstr = selectitemstr.substring(0, selectitemstr.length - 1);
		return selectitemstr;
	};

	$scope.Openskdd = function(id) {
		$state.go("app.skdd", {
			id: id
		}, {
			reload: true
		});
	};
	$scope.searchpost = function() {
		$scope.backPage.PageIndex=1;
		var url = config.HttpUrl + "/action/getcurriculums";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": $scope.form.Begindatestr,
			"Enddate": $scope.form.Endatestr,
			"State": -1,
			"Classroomid": null,
			"Teacherids": $scope.GetselectTeacherItem(),
			"Collegeids": $scope.GetselectcollegeItems(),
			"Majorids": $scope.GetselectmajorItems(),
			"Classesids": $scope.GetselectclasseItems(),
			"Curriculumsids": "",
			"Searhtxt": $scope.form.Search,
			"Campusid": Number($scope.Campus.selectCampusitem.Campusid),
			"Buildingid": Number($scope.building.selectbuildingitem.Buildingid),
			"Floorsid": Number($scope.floors.selectfloorsitem.Floorsid),
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.getcurriculumsList = data.Result.PageData;

				//   分页
				var objPage = {
					PageCount: data.Result.PageCount,
					PageIndex: data.Result.PageIndex,
					PageSize: data.Result.PageSize,
					RecordCount: data.Result.PageCount
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				$scope.backPage = pageFn(objPage, 5);

				setTimeout(function() {
					$('#showtable').trigger('footable_redraw');
				}, 1000);
			} else {
				$scope.getcurriculumsList = [];

				//   分页
				var objPage = {
					PageCount: 0,
					PageIndex: $scope.backPage.PageIndex,
					PageSize: $scope.backPage.PageSize,
					RecordCount: 0
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				$scope.backPage = pageFn(objPage, 5);
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	/*  -------------------- 分页、页码  -----------------------  */
	/*----------------
	//    分页对象添加页码
	//    return  obj  分页对象
	//    pagedata:obj  分页对象
	//    maxpagenumber:int  显示页码数默认5个页码
	------------------*/
	var pageFn = function(pagedata, maxpagenumber) {
		if(pagedata.length < 1) return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber) {
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1) {
				Snumber = 1
			} else if((pagedata.PageIndex + mid) > pagedata.PageCount) {
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			} else {
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++) {
				nub.push(Snumber + i);
			}
		} else {
			for(var i = 0; i < pagedata.PageCount; i++) {
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	};
	//  翻页
	$scope.pageClick = function(pageindex) {
		if(!(Number(pageindex) > 0)) return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount) {
			$scope.backPage.PageIndex = pageindex;
			$scope.searchpost();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */

	var run = function(){
		var myDate = new Date();
		$scope.form.Begindatestr = $filter('date')(myDate, 'yyyy-MM-dd') + " 00:00:00";
		$scope.form.Endatestr = $filter('date')(myDate, 'yyyy-MM-dd') + " 23:59:59";
		init_data();
		$scope.searchpost();
	}
	run();


}]);
