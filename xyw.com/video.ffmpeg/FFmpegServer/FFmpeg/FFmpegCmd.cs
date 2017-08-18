using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading;

namespace FFmpegServer
{
    public class FFmpegCmd
    {
        Process process = null;
        Thread th = null;

        public bool Execute(string ffmpeg, string args)
        {
            bool ret = false;
            if (!Running() && !string.IsNullOrEmpty(ffmpeg))
            {
                process = new Process();
                try
                {
                    ProcessStartInfo startInfo = new ProcessStartInfo();
                    startInfo.FileName = "cmd.exe";             // 设定需要执行的命令  
                    startInfo.Arguments = string.Format(" /C start {0}", ffmpeg + args);      // /C表示执行完命令后马上退出  
                    startInfo.UseShellExecute = false;          // 不使用系统外壳程序启动  
                    startInfo.RedirectStandardInput = true;    // 不重定向输入  
                    startInfo.RedirectStandardOutput = true;    // 重定向输出  
                    startInfo.RedirectStandardError = true;
                    startInfo.CreateNoWindow = true;            // 不创建窗口  
                    process.StartInfo = startInfo;

                    Console.WriteLine(startInfo.FileName + startInfo.Arguments);
                    ret = process.Start();

                    process.StandardInput.AutoFlush = true;
                    this.th = Thread.CurrentThread;
                }
                catch (Exception ex)
                {
                    System.Console.WriteLine(ex.StackTrace);
                    #region DEBUG
#if DEBUG
                    throw ex;
#endif
                    #endregion
                }
            }

            return ret;
        }

        public bool Start(string ffmpeg, string args)
        {
            bool ret = false;
            if (!Running() && !string.IsNullOrEmpty(ffmpeg))
            {
                process = new Process();
                try
                {
                    ProcessStartInfo startInfo = new ProcessStartInfo();
                    startInfo.FileName = "cmd.exe ";             // 设定需要执行的命令  
                    startInfo.Arguments = string.Format("/C {0}", ffmpeg + args);      // /C表示执行完命令后马上退出  
                    startInfo.UseShellExecute = false;          // 不使用系统外壳程序启动  
                    startInfo.RedirectStandardInput = true;    // 不重定向输入  
                    startInfo.RedirectStandardOutput = true;    // 重定向输出  
                    startInfo.RedirectStandardError = true;
                    startInfo.CreateNoWindow = true;            // 不创建窗口  
                    process.StartInfo = startInfo;

                    ret = process.Start();
                    process.StandardInput.AutoFlush = true;
                    this.th = Thread.CurrentThread;

                    //process.WaitForExit();                    
                    //this.process.Dispose();
                    this.process = null;
                    this.th = null;
                }
                catch (Exception ex)
                {
                    System.Console.WriteLine(ex.StackTrace);
                }
            }

            return ret;
        }

        public bool Running()
        {
            return (this.process != null && !this.process.HasExited);
        }

        public void Stop()
        {
            try
            {
                if (Running())
                {
                    this.process.StandardInput.Write("{0}", "q");
                    Thread.Sleep(10 * 1000);
                    if (!this.process.HasExited) { this.process.Kill(); }
                    this.process.Dispose();
                }
            }
            catch (Exception ex)
            {
                System.Console.WriteLine(ex.StackTrace);
            }

            this.process = null;
            this.th = null;
        }

        public static string GetCameraCmdString()
        {
            string dir = @"D:\vod\" + DateTime.Now.ToString("yyyyMMdd");
            string file = DateTime.Now.ToString("HHmmss") + ".mp4";
            string path = System.IO.Path.Combine(dir, file);
            string user = "admin";
            string pass = "xywadmin";
            string ip = "192.168.0.151";
            int port = 554;
            string istream = string.Format("rtsp://{0}:{1}@{2}:{3}/cam/realmonitor?channel=1&subtype=0", user, pass, ip, port);
            string cmd = string.Format("start ffmpeg -re -i {0} -acodec copy -vcodec libx264 -pix_fmt yuv420p -f mp4 {1}", istream, path);
            return cmd;
        }

        public static string GetScreenCmdString()
        {
            string dir = @"D:\vod\" + DateTime.Now.ToString("yyyyMMdd");
            string file = DateTime.Now.ToString("HHmmss") + ".mp4";
            string path = System.IO.Path.Combine(dir, file);

            string cmd = string.Format("start ffmpeg -f gdigrab -i desktop -f flv {0}", path);
            return cmd;
        }
    }

    public class VideoArgs
    {
        public int Fps { get; set; }            // 帧率
        public string Format { get; set; }      // 格式
        public string Codec { get; set; }       // 编码器
        public string PixelFormat { get; set; }     // 像素格式
        public string WxH { get; set; }             // 分辨率
        public int BitRate { get; set; }            // 码流率


    }

    public class AudioArgs
    {
        public int Channels { get; set; }       // 声道数       
        public int BitDepth { get; set; }       // 采样位深度
        public string Codec { get; set; }       // 编码器
        public string Format { get; set; }      // 格式
    }

    public class CmdMessage
    {
        public string CmdID { get; set; }
        public string CmdType { get; set; }
        public string JsonText { get; set; }

        public CmdMessage()
        {
            this.CmdID = "";
            this.CmdType = "";
            this.JsonText = "";
        }
    }

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
