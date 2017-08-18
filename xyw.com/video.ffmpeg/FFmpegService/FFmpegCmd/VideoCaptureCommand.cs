using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FFmpegService
{
    public class VideoCaptureCommand
    {
        public string CmdID { get; set; }
        public int CmdType { get; set; }
        public int CmdState { get; set; }
        public int ClassroomID { get; set; }

        public string TargetIP { get; set; }
        public int TargetPort { get; set; }
        public string TargetUser { get; set; }
        public string TargetPass { get; set; }

        public string CurriculumName { get; set; }      // 课程ID
        public string TeacherName { get; set; }         // 教师ID
        public string CharpterName { get; set; }        // 章节ID
        public int CurriculumDuration { get; set; }     // 课堂时长(秒)

        public string VideoFile { get; set; }    // 视频文件名
        public int VideoDuration { get; set; }    // 录制时长

        public VideoCaptureCommand()
        {
            this.CmdID = "";
            this.TargetIP = "";
            this.TargetUser = "";
            this.TargetPass = "";
            this.CurriculumName = "";
            this.TeacherName = "";
            this.CharpterName = "";
            this.VideoFile = "";
        }
    }
}
