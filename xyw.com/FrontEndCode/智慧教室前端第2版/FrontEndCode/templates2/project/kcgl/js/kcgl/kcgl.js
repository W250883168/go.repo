'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理
 */

/*    课程管理     */
app.controller("kcglContr", ['$scope', '$state', function($scope, $state) {
	//$state.go("app.qxgl",false);
    console.log("课程管理")


}]);





/*   弹窗  - 添加- 添加      */
app.controller("kcglKcjhAddAddContr", ['$scope', '$state', '$modal','$modalInstance','items', function($scope, $state, $modal,$modalInstance,items) {
	//$state.go("app.qxgl",false);
	console.log("课程管理-课程计划-添加-添加")
	
	console.log(items)
	$scope.items = items;
	
	
	$scope.from = {
			//   章节信息
			"section": {},
			//   课堂信息
			"classRoom": []
		}
		//    章节信息
	$scope.from.section = {
			//   课程ID
			"courseId": "",
			//   课程名称
			"courseName": "",
			//    章节ID
			"chapterId": "",
			//    章节名称
			"chapterName": "",
			//    顺序
			"order": 50,
			//   图标
			"pic": null,
			//    章节详情
			"chapterDetails": ""
		}
		//    课堂信息
	$scope.from.classRoom = [
		/*{
			//   班级
			"classId":"",
			"className":"",
			"date":"",
			"time":"",
			"classRoomId":"",
			"classRoomName":"",
			"teacherId":"",
			"teacherName":"",
			"isVideo":"1",
			"isPlay":"1",
			"isPlayComment":"1",
			"isLive":"1",
			"isLiveComment":"1"
		}*/
	];
	$scope.form = $scope.items;
	
	//  上传图片
	var handleFileSelect = function(evt) {
		var file = evt.currentTarget.files[0];
		if(!file && !/image\/\w+/.test(file.type))return false;
		var reader = new FileReader();
		reader.onload = function(evt) {
			$scope.$apply(function($scope) {
				$scope.from.section.pic = evt.target.result;
			});
		};
		reader.readAsDataURL(file);
	};
	$(document).on('change', '#fileInput_add_add', handleFileSelect);
	
	//    清除图片
	$scope.closePic = function(){
		$scope.from.section.pic = "";
	}


	//    打开弹窗  添加章节，添加课表
	$scope.modalOpenAddAdd = function() {
		var modalInstance = $modal.open({
			templateUrl: '../project/kcgl/html/kcgl/modal/modal_add_add_add.html',
			controller: 'kcglKcjhAddAddAddContr',
			windowClass: 'm-modal-kcgl-add-add',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(form) {
			console.log(form)
			if(!form) {
				//$scope.deviceText = "";
			} else {
				$scope.from.classRoom.push(form);
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	var clickDetails = function(){
		//    打开弹窗  添加章节，添加课表
		$scope.modalOpenAddAdd = function() {
			var modalInstance = $modal.open({
				templateUrl: '../project/kcgl/html/kcgl/modal/modal_add_add_add.html',
				controller: 'kcglKcjhAddAddAddContr',
				windowClass: 'm-modal-kcgl-add-add',
				resolve: {
					items: function() {
						return $scope.items;
					}
				}
			});
	
			modalInstance.result.then(function(form) {
				console.log(form)
				if(!form) {
					//$scope.deviceText = "";
				} else {
					$scope.from.classRoom.push(form);
				}
			}, function() {
				//$log.info('Modal dismissed at: ' + new Date());
			});
		}
	}
	
	//   操作
	$scope.operationClick = function(str){
		switch(str){
			case "details":
			break;
			case "edit":
			break;
			case "remove":
			break;
		}
	}
	
	
	$scope.ok = function() {
		//console.log(item)
		//   还差验证
		$modalInstance.close($scope.from);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
	
	
	//  run
	var run = function(){
		$scope.form = $scope.items;
	}
	//run();

}]);

/*   弹窗  - 添加 - 添加课表- 设置课堂信息      */
app.controller("kcglKcjhAddAddAddContr", ['$scope', '$state', '$modal', '$modalInstance', function($scope, $state, $modal, $modalInstance) {
	//$state.go("app.qxgl",false);
	console.log("课程管理-课程计划-添加-添加-设置课堂信息")

	//  编辑
	$scope.from = {
			//   班级
			"classId": "",
			"className": "",
			"date": "",
			"time": "",
			"classRoomId": "",
			"classRoomName": "",
			"teacherId": "",
			"teacherName": "",
			"isVideo": "1",
			"isPlay": "1",
			"isPlayComment": "1",
			"isLive": "1",
			"isLiveComment": "1"
		}
		//   上课教室显示标题
	$scope.classRoomIdTitle = "";
	//   上课老师显示标题
	$scope.teacherIdTitle = "";

	//   时间表
	$scope.dateTable = config.dateTable;

	$scope.openDate = function() {
		jeDate({
			dateCell: "#add_date",
			format: "YYYY-MM-DD",
			isTime: true,
			minDate: "2015-12-31",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.from.date = val;
			},
			okfun: function(elem, val) {
				$scope.from.date = val;
			},
			clearfun:function(elem, val) {
				$scope.from.date = "";
			}
		});
	}

	//    打开弹窗  选择班级
	$scope.modalOpenClass = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_class.html',
			controller: 'modalGetClassCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem) {
				$scope.from.classId = "";
				$scope.from.className = "";
			} else {
				$scope.from.classId = selectedItem.Id;
				$scope.from.className = selectedItem.Classesname;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择教室
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
			console.log(selectedItem)
			if(!selectedItem) {
				$scope.from.classRoomId = "";
				$scope.classRoomIdTitle = "";
			} else {
				$scope.from.classRoomId = selectedItem.addId;
				$scope.from.classRoomName = selectedItem.add;
				$scope.classRoomIdTitle = selectedItem.add;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择老师
	$scope.modalOpenTeacher = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_teacher.html',
			controller: 'modalGetTeacherCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem) {
				$scope.from.teacherId = "";
				$scope.teacherIdTitle = "";
			} else {
				$scope.from.teacherId = selectedItem.Usersid;
				$scope.from.teacherName = selectedItem.Truename;
				$scope.teacherIdTitle = selectedItem.Truename;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	$scope.ok = function() {
		//console.log(item)
		//   还差验证
		$modalInstance.close($scope.from);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};

}]);