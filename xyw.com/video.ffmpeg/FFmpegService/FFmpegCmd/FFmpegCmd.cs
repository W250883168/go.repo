using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace FFmpegService
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
                    #region // DEBUG
#if DEBUG
                    throw ex;
#endif
                    #endregion
                }
            }

            return ret;
        }

        public Process Execute2(string ffmpeg, string args)
        {
            System.Console.WriteLine(ffmpeg + args);
            process = new Process();
            process.StartInfo.UseShellExecute = false;          // 不使用系统外壳程序启动  
            process.StartInfo.RedirectStandardInput = true;    // 不重定向输入  
            process.StartInfo.RedirectStandardOutput = true;    // 重定向输出  
            process.StartInfo.RedirectStandardError = true;
            process.StartInfo.CreateNoWindow = true;            // 不创建窗口  
            process.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
            process.StartInfo.FileName = ffmpeg;
            process.StartInfo.Arguments = args;
            process.Start();

            new Thread(() =>
            {
                string date = DateTime.Now.ToString("yyyyMMdd");
                using (System.IO.StreamWriter tssw = new System.IO.StreamWriter(date + ".log", true))
                {
                    tssw.WriteLine(process.StandardError.ReadToEnd());
                }
            }).Start();


            return process;
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

        /// <summary>
        /// 运行cmd命令，不显示命令窗口
        /// </summary>
        /// <param name="cmdExe">应用程序完整的路径</param>
        /// <param name="cmdStr">执行命令的参数</param>
        /// <returns></returns>
        public Process RunProgram2(string cmdExe, string cmddir, string cmdStr)
        {
            Process proc = new Process();
            try
            {
                ProcessStartInfo info = new ProcessStartInfo();
                info.FileName = cmdExe;
                info.RedirectStandardInput = true;
                info.RedirectStandardOutput = true;
                info.RedirectStandardError = true;
                info.Arguments = cmdStr;
                // info.WindowStyle = ProcessWindowStyle.Hidden;
                // info.CreateNoWindow = true;
                info.UseShellExecute = false;
                proc.StartInfo = info;
                proc.Start();
            }
            catch (Exception ex)
            {
                System.Console.WriteLine(ex.ToString());
            }

            return proc;
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

                    Thread.Sleep(1 * 1000);
                    this.process.Dispose();
                    if (!this.process.HasExited) { this.process.Kill(); }
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






}
