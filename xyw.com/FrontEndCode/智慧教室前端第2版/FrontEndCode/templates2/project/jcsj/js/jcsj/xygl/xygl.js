'use strict';


app.controller("jcsjXyglContr", ['$scope', '$state','httpService', '$modal','toaster',function($scope, $state,httpService,$modal,toaster) {
	//$state.go("app.qxgl",false);
    console.log("学院管理")

    $scope.xyglData = {
    	"campusList":[],
    	"collegeList":[],
    	"majorList":[]
    }
    //    学院管理 树
    $scope.xygl_data = [];
    //    学院管理 右边列表
    $scope.xyglRightData = {
    	"item":{},
    	"list":[]
    };

    //   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }
	//   树 已展开 保存数组
	var treeIsOpens = [];

    //    getall
    var getAll = function(){
		//   是展开时加入展开数组
		treeIsOpens = [];
		dg_tree( $scope.xygl_data,function(item){
			if(item.expanded == true){
				switch(item.level){
					case 0:
						//
					break;
					case 1:
						treeIsOpens.push({"level":item.level,"Campusid":item.Id,"expanded":true});
					break;
					case 2:
						treeIsOpens.push({"level":item.level,"Collegeid":item.Id,"expanded":true});
					break;
					case 3:
						treeIsOpens.push({"level":item.level,"Majorid":item.Majorid,"expanded":true});
					break;
				}
			}
		});
		//
        var url = config.HttpUrl+"/basicset/getall";
        var data={};
        var promise =httpService.ajaxGet(url,null);
        promise.then(function (data) {
        	console.log("getall",data)
            if(data.Rcode=="1000"){
            	$scope.xyglData.campusList = data.Result[3];
            	$scope.xyglData.collegeList = data.Result[4];
            	$scope.xyglData.majorList = data.Result[5];

				$scope.xygl_data = outTree($scope.xyglData.campusList,$scope.xyglData.collegeList,$scope.xyglData.majorList);
				$scope.xygl_data[0].expanded = true;

				var objitem =$scope.xyglRightData.item;
				//	默认展开
				if(objitem == null || angular.equals({},objitem)){
					$scope.xygl_data[0].selected = true;
					$scope.xyglRightData.item = $scope.xygl_data[0];
					collegelist();
				}
				//  加入选中
				dg_tree( $scope.xygl_data,function(item){
					if($scope.xyglRightData.item.level == item.level && $scope.xyglRightData.item.Id == item.Id){
						item.selected = true;
					}else{
						item.selected = false;
					}
				});
				//   加入展开
				for(var c = 0; c < treeIsOpens.length; c++){
					if(treeIsOpens[c].level == 1){
						dg_tree( $scope.xygl_data,function(item){
							if(item.level == 1 && item.Id == treeIsOpens[c].Campusid){
								item.expanded = true;
								//return;
							}
						});
					}
					if(treeIsOpens[c].level == 2){
						dg_tree( $scope.xygl_data,function(item){
							if(item.level == 2 && item.Id == treeIsOpens[c].Collegeid){
								item.expanded = true;
								//return;
							}
						});
					}
					if(treeIsOpens[c].level == 3){
						dg_tree( $scope.xygl_data,function(item){
							if(item.level == 3 && item.Majorid == treeIsOpens[c].Majorid){
								item.expanded = true;
								//return;
							}
						});
					}
				}
            }else if(data.Rcode=="1002"){
            	$scope.xygl_data = [{label:'全部学院',children:[],level:0}];
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };

	//  递归树
	var dg_tree = function(tree, fn) {
		for(var i = 0; i < tree.length; i++) {
			fn(tree[i]);
			//  console.log(tree[i],tree[i].children)
			if(tree[i].children.length > 0) {
				dg_tree(tree[i].children, fn);
			}
		}
	}

	/*    生成学院学科教室树      */
   //   obj,obj,obj
	var outTree = function(campusList,collegeList,majorList){
		var tree = [];
		//  计数
		var n1 = 0,n2 = 0;
		//   学院
	   	for(var c in campusList){
	   		//  名称
	   		campusList[c].label = campusList[c].Collegename;
	   		campusList[c].level = 1;
	   		tree.push(campusList[c]);
	   		n2 = 0;
	   		//   学科
	   		for(var b in collegeList){
				if(!("children" in tree[n1])) tree[n1].children = [];
	   			if(tree[n1].Id == collegeList[b].Collegeid){
	   				//  名称
	   				collegeList[b].label = collegeList[b].Majorname;
	   				collegeList[b].level = 2;
	   				//   插入数组
	   				tree[n1].children.push(collegeList[b]);
					// 班级
					for(var a in majorList){
						if(!("children" in tree[n1].children[n2])) tree[n1].children[n2].children=[];
						//    加上children属性
						if(!("children" in majorList[a])) majorList[a].children = [];
						if(tree[n1].children[n2].Id == majorList[a].Majorid){
							majorList[a].label = majorList[a].Classesname;
							majorList[a].level = 3;
							tree[n1].children[n2].children.push(majorList[a]);
						}
					}
					n2++;
	   			}
	   		}
	   		n1++;
	   	}
	   	return [{label:'全部学院',children:tree,level:0}];
	}

	//    学院查询
    var collegelist = function(){
        var url = config.HttpUrl+"/system/bs/collegelist";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"PageIndex": -1
			//"PageSize": $scope.backPage.PageSize
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("学院查询",data)
            if(data.Rcode=="1000"){
            	$scope.xyglRightData.list = data.Result.PageData;
            	//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
            }else{
                $scope.xyglRightData.list = [];
            }
        }, function (reason) {}, function (update) {});
    };

    //    删除学院
    var collegedel = function(Id){
    	if(!Id)return false;
        var url = config.HttpUrl+"/system/bs/collegedel";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			      "Id":Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("删除学院",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '删除成功！');
            collegelist();
            getAll();
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };
    //    科系查询
    var majorlist = function(){
        var url = config.HttpUrl+"/system/bs/majorlist";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Collegeid":$scope.xyglRightData.item.Id,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};

        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("科系查询",data)
            if(data.Rcode=="1000"){
            	if(data.Result != null){
            		$scope.xyglRightData.list = data.Result.PageData;
            		//   分页
					var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
					if((objPage.RecordCount % objPage.PageSize)==0){
						objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
					}else{
						objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
					}
					$scope.backPage = pageFn(objPage,5);
            	}else{
            		$scope.xyglRightData.list = [];
            	}
            }else{
                console.log(data.Reason);
                $scope.xyglRightData.list = [];
            }
        }, function (reason) {}, function (update) {});
    };
    //删除科系
    var majordel = function(Id){
    	if(!Id)return false;
        var url = config.HttpUrl+"/system/bs/majordel";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("删除科系",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '删除成功！');
            majorlist();
            getAll();
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };
     //    班级查询
    var classeslist = function(){
        var url = config.HttpUrl+"/system/bs/classeslist";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Majorid":$scope.xyglRightData.item.Id,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("班级查询",data)
            if(data.Rcode=="1000"){
            	if(data.Result != null){
            		$scope.xyglRightData.list = data.Result.PageData;
            		//   分页
					var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
					if((objPage.RecordCount % objPage.PageSize)==0){
						objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
					}else{
						objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
					}
					$scope.backPage = pageFn(objPage,5);
            	}else{
            		$scope.xyglRightData.list = [];
            	}
            }else{
                $scope.xyglRightData.list = [];
            }
        }, function (reason) {}, function (update) {});
    };

    //删除班级
    var classesdel = function(Id){
        var url = config.HttpUrl+"/system/bs/classesdel";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("删除班级",data)
            if(data.Rcode=="1000"){
            	classeslist();
              getAll();
              toaster.pop('success', '删除成功！');
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };
    //添加按钮功能
    $scope.openModalAdd = function (str,active,item) {
    	if(!("level" in active)){
        $modal.open({
          templateUrl: 'modal/modal_alert_all.html',
          controller: 'modalAlert2Conter',
          resolve: {
            items: function () {
              return {"type":'info',"msg":'请选择选项'};
            }
          }
        });
        return false;
      }
        var modalInstance = $modal.open({
            templateUrl: '../project/jcsj/html/jcsj/xygl/modal_add.html',
            controller: 'modalJcsjXyglContr',
            windowClass: 'm-modal-xygl',
            resolve: {
                items: function () {
                    return {'str':str,"active":active,"item":item};
                }
            }
        });
      modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul){
				//   刷新列表
				switch($scope.xyglRightData.item.level){
		    		case 0:
		    			collegelist();
		    		break;
		    		case 1:
		    			majorlist();
		    		break;
		    		case 2:
		    			classeslist();
		    		break;
					case 3:

					break;

		    	}
				getAll();
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
    }

    //    树点击
    $scope.xygl_tree_handler = function(branch) {
		console.log(branch);
    	$scope.backPage.PageIndex = 1;
    	$scope.backPage.PageSize = 10;
    	$scope.xyglRightData.item = branch;
    	switch($scope.xyglRightData.item.level){
    		case 0:
    			collegelist();
    		break;
    		case 1:
    			majorlist();
    		break;
    		case 2:
    			classeslist();
    		break;
			case 3:
				$scope.xyglRightData.list = [];
			break;
    	}
		if(branch.level != 0){
			$scope.xygl_data[0].selected = false;
		}
    };

    //   删除弹窗
    $scope.deleteItem = function(active,item){
      var modalInstance = $modal.open({
        templateUrl: 'modal/modal_alert_all.html',
        controller: 'modalAlert2Conter',
        resolve: {
          items: function () {
            return {"type":'warning',"msg":'你确定要删除吗？'};
          }
        }
      });
      modalInstance.result.then(function(bul){
        if(bul){
          switch(active.level){
            case 0:
              //   删除学院
              collegedel(item.Id);
              break;
            case 1:
              majordel(item.Id);
              break;
            case 2:
              classesdel(item.Id);
              break;
          }
        }
      });
    }

    /*  -------------------- 分页、页码  -----------------------  */
	//$scope.backPage = {};
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
	};
	//  翻页
	$scope.pageClick = function(pageindex){
		if(!(Number(pageindex) > 0))return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount){
			$scope.backPage.PageIndex = pageindex;
			switch($scope.xyglRightData.item.level){
	    		case 0:
    				collegelist();
	    		break;
	    		case 1:
	    			majorlist();
	    		break;
	    		case 2:
	    			classeslist();
	    		break;
	    	}
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */
	$scope.run = function(){
		getAll();
	}
	$scope.run();

}]);


app.controller("modalJcsjXyglContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("基础数据-学院管理-弹窗");
	$scope.items = items;

	console.log(items);

	$scope.form = {
		//   学院
		"college":{},
		//   科系
		"major":{},
		//   班级
		"classes":{},

	}
	$scope.form.college = {
		"Collegecode":"",
        "Collegename":"",
        "Collegeicon":"",
        "Campusid":1,
        "CampusItems":[],
        "Id":null
	}

	$scope.form.major = {
		"Majorcode":"",
        "Majorname":"",
        "Majoricon":"",
        "Collegeid":null,
        "Id":null
	}
	$scope.form.classes = {
		"Classescode":"",
        "Classesname":"",
        "Classesicon":"",
        "Majorid":null,
        "Id":null,
		"Enrollmentyear":null
	}

	//   上传图片
	$scope.upimglist = [];

	//    清除图片
	$scope.closePic = function(index){
		$scope.upimglist.splice(index,1);
	}
	
	
	//    验证学院代码
    $scope.changeSubjectcode = function(){
    	//   本级学院增加长度
    	var xk_lang = 4;
    	//   学院代码
    	var code = '';
    	//   
    	if($scope.items.active.level == 0){
    		code = $scope.form.college.Collegecode;
    	}else if($scope.items.active.level == 1){
    		code = $scope.form.major.Majorcode;
    	}else if($scope.items.active.level == 2){
    		code = $scope.form.classes.Classescode;
    	}
    	
    	//   上级学院代码
    	var pcode = "";
    	if($scope.items.active.level == 0){
    		pcode = "";
    	}else if($scope.items.active.level == 1){
    		pcode = $scope.items.active.Collegecode;
    	}else if($scope.items.active.level == 2){
    		pcode = $scope.items.active.Majorcode;
    	}
    	
    	//    验证输入必须为数字
    	if(!(/^[a-zA-Z0-9-]*$/.test(code))){
    		
    		$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'只能输入上级代码加4位的数字或字母！'};
                }
              }
            });
    		//alert('只能输入上级代码加4位的数字！');
    		code = pcode;
    	}

    	//
    	if(code.length < pcode.length){
			//   里面没有上级学院代码  则 放入上级学院代码
			code = pcode;
    	}else{
    		//   有  验证前面的数字是不是上级学院代码
    		if(code.substr(0,pcode.length) == pcode){
    			//   是上级学院代码 。限制输入长度
    			if(code.length > pcode.length + xk_lang){
    				code = code.substr(0,code.length - 1);
    			}else{
    				console.log(code);
    			}
    		}else{
    			//   不是上级学院代码 。 放入上级学院代码
    			code = pcode;
    		}
    	}
    	//
    	//$scope.form.subject.Subjectcode = code;
    	if($scope.items.active.level == 0){
    		$scope.form.college.Collegecode = code;
    	}else if($scope.items.active.level == 1){
    		$scope.form.major.Majorcode = code;
    	}else if($scope.items.active.level == 2){
    		$scope.form.classes.Classescode = code;
    	}
    }
	


	//   取校区
	//    校区查询
    var campuslist = function(){
        var url = config.HttpUrl+"/system/bs/campuslist";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"PageIndex": -1
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("校区查询",data)
            if(data.Rcode=="1000"){
            	$scope.form.college.CampusItems = data.Result.PageData;
            }else{
                console.log(data.Reason);
                $scope.form.college.CampusItems = [];
            }
        }, function (reason) {}, function (update) {});
    };
    
    
    //    验证提交
    var yz_code = function(){
    	//
		var temp_bol = false;
		switch($scope.items.active.level){
			case 0:
				if( $scope.form.college.Collegecode == "" ||  $scope.form.college.Collegecode == undefined){
					temp_bol = true;
				}
			break;
			case 1:
				if($scope.form.major.Majorcode == $scope.items.active.Collegecode){
					temp_bol = true;	
				}
			break;
			case 2:
				if($scope.form.classes.Classescode == $scope.items.active.Majorcode){
					temp_bol = true;
				}
			break;
		}
		if(temp_bol){
    		$modal.open({
              templateUrl: 'modal/modal_alert_all.html',
              controller: 'modalAlert2Conter',
              resolve: {
                items: function () {
                  return {"type":'danger',"msg":'代码不能与上级代码相同！'};
                }
              }
            });
            return false;
		}else{
			return true;
		}
		
		
    }


	//   学院添加
    var collegeadd = function(){
		if(!(formValidate($scope.form.college.Collegecode).minLength(0).outMsg(2811).isOk))return false;
		if(!(formValidate($scope.form.college.Collegename).minLength(0).outMsg(2812).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
        var url = config.HttpUrl+"/system/bs/collegeadd";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Collegecode":$scope.form.college.Collegecode,
            "Collegename":$scope.form.college.Collegename,
            "Collegeicon":$scope.form.college.Collegeicon,
            "Campusid":Number($scope.form.college.Campusid)
			//"Id":0
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("学院添加",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '添加成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };

    //   学院修改
    var collegechange = function(){
		if(!(formValidate($scope.form.college.Collegecode).minLength(0).outMsg(2811).isOk))return false;
		if(!(formValidate($scope.form.college.Collegename).minLength(0).outMsg(2812).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
        var url = config.HttpUrl+"/system/bs/collegechange";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Collegecode":$scope.form.college.Collegecode,
            "Collegename":$scope.form.college.Collegename,
            "Collegeicon":$scope.form.college.Collegeicon,
            "Id":$scope.items.item.Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("学院修改",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '修改成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };
    /////////////////////////////
    //   科系添加
    var majoradd = function(){
		if(!(formValidate($scope.form.major.Majorcode).minLength(0).outMsg(2813).isOk))return false;
		if(!(formValidate($scope.form.major.Majorname).minLength(0).outMsg(2814).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
        var url = config.HttpUrl+"/system/bs/majoradd";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Majorcode":$scope.form.major.Majorcode,
            "Majorname":$scope.form.major.Majorname,
            "Majoricon":$scope.form.major.Majoricon,
            "Collegeid":$scope.form.major.Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("科系添加",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '添加成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };

    //   科系修改
    var majorchange = function(){
		if(!(formValidate($scope.form.major.Majorcode).minLength(0).outMsg(2813).isOk))return false;
		if(!(formValidate($scope.form.major.Majorname).minLength(0).outMsg(2814).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
        var url = config.HttpUrl+"/system/bs/majorchange";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Majorcode":$scope.form.major.Majorcode,
            "Majorname":$scope.form.major.Majorname,
            "Majoricon":$scope.form.major.Majoricon,
            "Collegeid":$scope.form.major.Collegeid,
            "Id":$scope.items.item.Id
		};

        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("科系修改",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '修改成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };
    /////////////////////////////////////////////

    //   班级添加
    var classesadd = function(){
		if(!(formValidate($scope.form.classes.Classescode).minLength(0).outMsg(2815).isOk))return false;
		if(!(formValidate($scope.form.classes.Classesname).minLength(0).outMsg(2816).isOk))return false;
		if(!(formValidate($scope.form.classes.Enrollmentyear).isNumber().outMsg(2817).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
        var url = config.HttpUrl+"/system/bs/classesadd";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Classescode":$scope.form.classes.Classescode,
            "Classesname":$scope.form.classes.Classesname,
            "Classesicon":$scope.form.classes.Classesicon,
            "Majorid":$scope.items.active.Id,
			"Enrollmentyear":Number($scope.form.classes.Enrollmentyear)
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("班级添加",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '添加成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };

    //   班级修改
    var classeschange = function(){
		if(!(formValidate($scope.form.classes.Classescode).minLength(0).outMsg(2815).isOk))return false;
		if(!(formValidate($scope.form.classes.Classesname).minLength(0).outMsg(2816).isOk))return false;
		if(!(formValidate($scope.form.classes.Enrollmentyear).isNumber().outMsg(2817).isOk))return false;
		
		//   验证提交
		if(!yz_code())return false;
		//
		$scope.changeSubjectcode();
		
		var url = config.HttpUrl+"/system/bs/classeschange";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Classescode":$scope.form.classes.Classescode,
            "Classesname":$scope.form.classes.Classesname,
            "Classesicon":$scope.form.classes.Classesicon,
            "Majorid":$scope.form.classes.Majorid,
            "Id":$scope.items.item.Id,
			"Enrollmentyear":Number($scope.form.classes.Enrollmentyear)
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("班级修改",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '修改成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };
    /////////////////////////////////////////////////////////////////
    //
	//    打开弹出 -选择日期
	$scope.showDate = function() {
		jeDate({
			dateCell: "#jd_begindate",
			format: "YYYY",
			isTime: true,
			minDate: "2000",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.classes.Enrollmentyear = val;
			},
			okfun: function(elem,val) {
				$scope.form.classes.Enrollmentyear = val;
			},
			clearfun:function(elem, val) {
				$scope.form.classes.Enrollmentyear = null;
			}
		});
	}
	/////////////////////////////////////////////////////////////////

	$scope.ok = function(){
		if($scope.items.str == "add"){
			switch($scope.items.active.level){
				case 0:
					if($scope.upimglist.length > 0){
						$scope.form.college.Collegeicon = $scope.upimglist[0].Result;
					}else{
						$scope.form.college.Collegeicon = "";
					}
					$scope.form.college.Id = 0;
					collegeadd();
				break;
				case 1:
					if($scope.upimglist.length > 0){
						$scope.form.major.Majoricon = $scope.upimglist[0].Result;
					}else{
						$scope.form.major.Majoricon = "";
					}
					$scope.form.major.Id = $scope.items.active.Id;
					majoradd();
				break;
				case 2:
					if($scope.upimglist.length > 0){
						$scope.form.classes.Classesicon = $scope.upimglist[0].Result;
					}else{
						$scope.form.classes.Classesicon = "";
					}
					$scope.form.classes.Majorid = $scope.items.active.Id;
					classesadd();
				break;
			}
		}
		if($scope.items.str == "edit"){
			switch($scope.items.active.level){
				case 0:
					if($scope.upimglist.length > 0){
						$scope.form.college.Collegeicon = $scope.upimglist[0].Result;
					}else{
						$scope.form.college.Collegeicon = "";
					}
					collegechange();
				break;
				case 1:
					if($scope.upimglist.length > 0){
						$scope.form.major.Majoricon = $scope.upimglist[0].Result;
					}else{
						$scope.form.major.Majoricon = "";
					}
//					$scope.form.major.Id = $scope.items.active.Id;
					majorchange();
				break;
				case 2:
					if($scope.upimglist.length > 0){
						$scope.form.classes.Classesicon = $scope.upimglist[0].Result;
					}else{
						$scope.form.classes.Classesicon = "";
					}
					//$scope.form.classes.Majorid = $scope.items.active.Majorid;
					classeschange();
				break;
			}
		}
	}
	//  选择校区时传ID
	$scope.changeCampusItem = function (item) {
    $scope.form.college.Campusid = item.Campusid;
  }
	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	$scope.run = function(){
		campuslist();

		if($scope.items.str == "add"){
			switch($scope.items.active.level){
				case 0:
					$scope.form.college.Collegecode = "";
				break;
				case 1:
					$scope.form.major.Majorcode = $scope.items.active.Collegecode;
				break;
				case 2:
					$scope.form.classes.Classescode = $scope.items.active.Majorcode;
				break;
			}
		}
		if($scope.items.str == "edit"){
			
			switch($scope.items.active.level){
				case 0:
					$scope.form.college.Collegecode = $scope.items.item.Collegecode;
					$scope.form.college.Collegename = $scope.items.item.Collegename;
					$scope.form.college.Collegeicon = $scope.items.item.Collegeicon;
					$scope.form.college.Id = 0;
					if($scope.form.college.Collegeicon.length > 0)$scope.upimglist[0] = {Result:$scope.form.college.Collegeicon};
				break;
				case 1:
		            $scope.form.major = $.extend({}, $scope.form.major, $scope.items.item);
		            $scope.form.major.Collegeid = $scope.items.active.Id;
					if($scope.form.major.Majoricon.length > 0)$scope.upimglist[0] = {Result:$scope.form.major.Majoricon};
				break;
				case 2:
					$scope.form.classes = $.extend({}, $scope.form.classes, $scope.items.item);
		            $scope.form.classes.Id = $scope.items.active.Id;
					if($scope.form.classes.Classesicon.length > 0)$scope.upimglist[0] = {Result:$scope.form.classes.Classesicon};
					//   入学年份
					//if($scope.items.item.Enrollmentyear == 0)$scope.items.item.Enrollmentyear ='';
					//$scope.form.classes.Enrollmentyear = 
				break;
			}
		}
	}
	$scope.run();
}]);

