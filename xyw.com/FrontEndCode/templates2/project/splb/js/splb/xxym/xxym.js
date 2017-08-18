//视频录播页面
app.controller("splbXxymindexContr",['$scope','httpService','$modal','$location','$state','formValidate','toaster',function($scope,httpService,$modal,$location,$state,formValidate,toaster){
	console.log("视频录播-详细页面",$location.search().operation);
	//  视频ID
  $scope.videoId = $location.search().vid;
  //  中间ID
  $scope.CurriculumClassroomChapterID = $location.search().cid;
  //  操作状态
  $scope.operation = $location.search().operation;
	//   上传图片
	$scope.upimglist = [];

	//    清除图片
	$scope.closePic = function(index){
		$scope.upimglist.splice(index,1);
	}

	//  form
	$scope.form = {
		"ID":null,
		"VideoTitle":"",
		"VideoInfo":"",
		"CoverImage":"",
		"LiveState":null,
		"BeginTime":"",
		"EndTime":"",
		"VideoDuration":null,
		"RecommendReads":"",
		"VodPath1":"",
		"VodPath2":"",
		"AllowVOD":1,
		// ---
		"AllowVODItem":{"val":1,"title":"是"},
		//  ---
		"AllowVODItems":[
			{"val":1,"title":"是"},
			{"val":0,"title":"否"}
		],
		"IsRelease":0,
		"AllowComments":1,
		"IsCheckComments":1,
		"AllowDownload":1,
		"PlayNum":0,
		"DownloadNum":0,
		"AttachmentList":[]
	}

	//
	$scope.videoList = [];


	//
	$scope.changeAllowVODItem = function(item){
		$scope.form.AllowVODItem = item;
		$scope.form.AllowVOD = item.val;
	}


	/**
	 * 取视频详情
	 */
	var videodetails = function(id) {
		var url = config.HttpUrl + "/vod/videodetails";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "ID":Number(id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取视频详情",data)
			if(data.Rcode == "1000") {
				$scope.form = $.extend({},$scope.form,data.Result);
				//    载入到图片
				if($scope.form.CoverImage){
					$scope.upimglist = [{"Result":$scope.form.CoverImage}];
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}


	/**
	 * 视频详情修改
	 */
	var updatevideodetails = function(id) {


		var url = config.HttpUrl + "/vod/updatevideodetails";
		var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "ID":$scope.form.ID,
			"VideoTitle":$scope.form.VideoTitle,
			"VideoInfo":$scope.form.VideoInfo,
			"CoverImage":$scope.form.CoverImage,
			"RecommendReads":$scope.form.RecommendReads,
			"DownloadNum":Number($scope.form.DownloadNum),
			"PlayNum":Number($scope.form.PlayNum),
			"AllowDownload":Number($scope.form.AllowDownload),
			"AllowComments":Number($scope.form.AllowComments),
			"IsCheckComments":Number($scope.form.IsCheckComments),
			"IsRelease":Number($scope.form.IsRelease),
			"LiveState":Number($scope.form.LiveState)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("视频详情修改",data)
			if(data.Rcode == "1000") {
        toaster.pop('success','修改成功');
        $state.go("^",{},{reload:true});
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	$scope.openModalAddSp = function (str) {
        var modalInstance = $modal.open({
            templateUrl: '../project/splb/html/splb/modal_sp.html',
            controller: 'modalSpContr',
            windowClass: 'm-modal-splb',
            resolve: {
                items: function () {
                    return {"item":$scope.form,"videoUrl":str};
                }
            }
        });
    }



	$scope.cancel=function(){
		$state.go("^");
	}

	$scope.ok = function(){
		if($scope.upimglist.length > 0){
			//   取最后一个图片
			$scope.form.CoverImage = $scope.upimglist[$scope.upimglist.length -1].Result;
		}

		updatevideodetails();
	}

	$scope.run = function(){
		videodetails($scope.videoId);
	}
	$scope.run();



}]);

//视频录播-视频播放弹窗
app.controller("modalSpContr",['$scope', 'httpService', '$modal', '$modalInstance','items',function ($scope, httpService, $modal,$modalInstance,items) {
	console.log("视频录播-视频播放弹窗");
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	//
	$scope.items = items;

}]);
