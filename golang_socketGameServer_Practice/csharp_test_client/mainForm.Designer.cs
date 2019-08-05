using System;

namespace csharp_test_client
{
    partial class mainForm
    {
        /// <summary>
        /// 필수 디자이너 변수입니다.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// 사용 중인 모든 리소스를 정리합니다.
        /// </summary>
        /// <param name="disposing">관리되는 리소스를 삭제해야 하면 true이고, 그렇지 않으면 false입니다.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form 디자이너에서 생성한 코드

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.components = new System.ComponentModel.Container();
            this.btnDisconnect = new System.Windows.Forms.Button();
            this.btnConnect = new System.Windows.Forms.Button();
            this.groupBox5 = new System.Windows.Forms.GroupBox();
            this.textBoxPort = new System.Windows.Forms.TextBox();
            this.label10 = new System.Windows.Forms.Label();
            this.checkBoxLocalHostIP = new System.Windows.Forms.CheckBox();
            this.textBoxIP = new System.Windows.Forms.TextBox();
            this.label9 = new System.Windows.Forms.Label();
            this.labelStatus = new System.Windows.Forms.Label();
            this.listBoxLog = new System.Windows.Forms.ListBox();
            this.label1 = new System.Windows.Forms.Label();
            this.textBoxUserID = new System.Windows.Forms.TextBox();
            this.textBoxUserPW = new System.Windows.Forms.TextBox();
            this.label2 = new System.Windows.Forms.Label();
            this.button2 = new System.Windows.Forms.Button();
            this.Room = new System.Windows.Forms.GroupBox();
            this.labelMasterUser = new System.Windows.Forms.Label();
            this.label5 = new System.Windows.Forms.Label();
            this.btnRoomGameStart = new System.Windows.Forms.Button();
            this.btnRoomBattingTie = new System.Windows.Forms.Button();
            this.btnRoomBattingPlayer = new System.Windows.Forms.Button();
            this.btnRoomBattingBanker = new System.Windows.Forms.Button();
            this.btnRoomWhisper = new System.Windows.Forms.Button();
            this.textBoxRoomSendWhisperRecieverId = new System.Windows.Forms.TextBox();
            this.textBoxRoomSendWhisperMsg = new System.Windows.Forms.TextBox();
            this.btnRoomChat = new System.Windows.Forms.Button();
            this.textBoxRoomSendMsg = new System.Windows.Forms.TextBox();
            this.listBoxRoomChatMsg = new System.Windows.Forms.ListBox();
            this.label4 = new System.Windows.Forms.Label();
            this.listBoxRoomUserList = new System.Windows.Forms.ListBox();
            this.btn_RoomLeave = new System.Windows.Forms.Button();
            this.btn_RoomEnter = new System.Windows.Forms.Button();
            this.textBoxRoomNumber = new System.Windows.Forms.TextBox();
            this.label3 = new System.Windows.Forms.Label();
            this.digitalTime = new System.Windows.Forms.Label();
            this.timer1 = new System.Windows.Forms.Timer(this.components);
            this.groupBox5.SuspendLayout();
            this.Room.SuspendLayout();
            this.SuspendLayout();
            this.btnDisconnect.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnDisconnect.Location = new System.Drawing.Point(895, 77);
            this.btnDisconnect.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnDisconnect.Name = "btnDisconnect";
            this.btnDisconnect.Size = new System.Drawing.Size(125, 47);
            this.btnDisconnect.TabIndex = 29;
            this.btnDisconnect.Text = "접속 끊기";
            this.btnDisconnect.UseVisualStyleBackColor = true;
            this.btnDisconnect.Click += new System.EventHandler(this.btnDisconnect_Click);
            this.btnConnect.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnConnect.Location = new System.Drawing.Point(894, 27);
            this.btnConnect.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnConnect.Name = "btnConnect";
            this.btnConnect.Size = new System.Drawing.Size(126, 50);
            this.btnConnect.TabIndex = 28;
            this.btnConnect.Text = "접속하기";
            this.btnConnect.UseVisualStyleBackColor = true;
            this.btnConnect.Click += new System.EventHandler(this.btnConnect_Click);
            this.groupBox5.Controls.Add(this.textBoxPort);
            this.groupBox5.Controls.Add(this.label10);
            this.groupBox5.Controls.Add(this.checkBoxLocalHostIP);
            this.groupBox5.Controls.Add(this.textBoxIP);
            this.groupBox5.Controls.Add(this.label9);
            this.groupBox5.Controls.Add(this.btnConnect);
            this.groupBox5.Controls.Add(this.btnDisconnect);
            this.groupBox5.Location = new System.Drawing.Point(27, 25);
            this.groupBox5.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.groupBox5.Name = "groupBox5";
            this.groupBox5.Padding = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.groupBox5.Size = new System.Drawing.Size(1038, 139);
            this.groupBox5.TabIndex = 27;
            this.groupBox5.TabStop = false;
            this.groupBox5.Text = "Socket 더미 클라이언트 설정";
            this.textBoxPort.Location = new System.Drawing.Point(426, 64);
            this.textBoxPort.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxPort.MaxLength = 6;
            this.textBoxPort.Name = "textBoxPort";
            this.textBoxPort.Size = new System.Drawing.Size(72, 31);
            this.textBoxPort.TabIndex = 18;
            this.textBoxPort.Text = "11021";
            this.textBoxPort.WordWrap = false;
            this.label10.AutoSize = true;
            this.label10.Location = new System.Drawing.Point(324, 64);
            this.label10.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label10.Name = "label10";
            this.label10.Size = new System.Drawing.Size(94, 25);
            this.label10.TabIndex = 17;
            this.label10.Text = "포트 번호:";
            this.checkBoxLocalHostIP.AutoSize = true;
            this.checkBoxLocalHostIP.Checked = true;
            this.checkBoxLocalHostIP.CheckState = System.Windows.Forms.CheckState.Checked;
            this.checkBoxLocalHostIP.Location = new System.Drawing.Point(577, 66);
            this.checkBoxLocalHostIP.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.checkBoxLocalHostIP.Name = "checkBoxLocalHostIP";
            this.checkBoxLocalHostIP.Size = new System.Drawing.Size(152, 29);
            this.checkBoxLocalHostIP.TabIndex = 15;
            this.checkBoxLocalHostIP.Text = "localhost 사용";
            this.checkBoxLocalHostIP.UseVisualStyleBackColor = true;
            this.textBoxIP.Location = new System.Drawing.Point(159, 64);
            this.textBoxIP.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxIP.MaxLength = 6;
            this.textBoxIP.Name = "textBoxIP";
            this.textBoxIP.Size = new System.Drawing.Size(123, 31);
            this.textBoxIP.TabIndex = 11;
            this.textBoxIP.Text = "0.0.0.0";
            this.textBoxIP.WordWrap = false;
            this.label9.AutoSize = true;
            this.label9.Location = new System.Drawing.Point(57, 64);
            this.label9.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label9.Name = "label9";
            this.label9.Size = new System.Drawing.Size(94, 25);
            this.label9.TabIndex = 10;
            this.label9.Text = "서버 주소:";
            this.labelStatus.AutoSize = true;
            this.labelStatus.Location = new System.Drawing.Point(27, 1394);
            this.labelStatus.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.labelStatus.Name = "labelStatus";
            this.labelStatus.Size = new System.Drawing.Size(166, 25);
            this.labelStatus.TabIndex = 40;
            this.labelStatus.Text = "서버 접속 상태: ???";
            this.listBoxLog.FormattingEnabled = true;
            this.listBoxLog.HorizontalScrollbar = true;
            this.listBoxLog.ItemHeight = 25;
            this.listBoxLog.Location = new System.Drawing.Point(27, 1108);
            this.listBoxLog.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.listBoxLog.Name = "listBoxLog";
            this.listBoxLog.Size = new System.Drawing.Size(1038, 279);
            this.listBoxLog.TabIndex = 41;
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(436, 232);
            this.label1.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(71, 25);
            this.label1.TabIndex = 42;
            this.label1.Text = "UserID:";
            this.textBoxUserID.Location = new System.Drawing.Point(515, 226);
            this.textBoxUserID.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxUserID.MaxLength = 6;
            this.textBoxUserID.Name = "textBoxUserID";
            this.textBoxUserID.Size = new System.Drawing.Size(123, 31);
            this.textBoxUserID.TabIndex = 43;
            this.textBoxUserID.Text = "jacking75";
            this.textBoxUserID.WordWrap = false;
            this.textBoxUserPW.Location = new System.Drawing.Point(746, 225);
            this.textBoxUserPW.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxUserPW.MaxLength = 6;
            this.textBoxUserPW.Name = "textBoxUserPW";
            this.textBoxUserPW.Size = new System.Drawing.Size(123, 31);
            this.textBoxUserPW.TabIndex = 45;
            this.textBoxUserPW.Text = "jacking75";
            this.textBoxUserPW.WordWrap = false;
            this.label2.AutoSize = true;
            this.label2.Location = new System.Drawing.Point(657, 226);
            this.label2.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(81, 25);
            this.label2.TabIndex = 44;
            this.label2.Text = "PassWD:";
            this.button2.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.button2.Location = new System.Drawing.Point(905, 216);
            this.button2.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.button2.Name = "button2";
            this.button2.Size = new System.Drawing.Size(142, 49);
            this.button2.TabIndex = 46;
            this.button2.Text = "Login";
            this.button2.UseVisualStyleBackColor = true;
            this.button2.Click += new System.EventHandler(this.button2_Click);
            this.Room.Controls.Add(this.labelMasterUser);
            this.Room.Controls.Add(this.label5);
            this.Room.Controls.Add(this.btnRoomGameStart);
            this.Room.Controls.Add(this.btnRoomBattingTie);
            this.Room.Controls.Add(this.btnRoomBattingPlayer);
            this.Room.Controls.Add(this.btnRoomBattingBanker);
            this.Room.Controls.Add(this.btnRoomWhisper);
            this.Room.Controls.Add(this.textBoxRoomSendWhisperRecieverId);
            this.Room.Controls.Add(this.textBoxRoomSendWhisperMsg);
            this.Room.Controls.Add(this.btnRoomChat);
            this.Room.Controls.Add(this.textBoxRoomSendMsg);
            this.Room.Controls.Add(this.listBoxRoomChatMsg);
            this.Room.Controls.Add(this.label4);
            this.Room.Controls.Add(this.listBoxRoomUserList);
            this.Room.Controls.Add(this.btn_RoomLeave);
            this.Room.Controls.Add(this.btn_RoomEnter);
            this.Room.Controls.Add(this.textBoxRoomNumber);
            this.Room.Controls.Add(this.label3);
            this.Room.Location = new System.Drawing.Point(27, 298);
            this.Room.Margin = new System.Windows.Forms.Padding(4, 6, 4, 6);
            this.Room.Name = "Room";
            this.Room.Padding = new System.Windows.Forms.Padding(4, 6, 4, 6);
            this.Room.Size = new System.Drawing.Size(1038, 797);
            this.Room.TabIndex = 47;
            this.Room.TabStop = false;
            this.Room.Text = "Room";
            this.labelMasterUser.AutoSize = true;
            this.labelMasterUser.Location = new System.Drawing.Point(309, 80);
            this.labelMasterUser.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.labelMasterUser.Name = "labelMasterUser";
            this.labelMasterUser.Size = new System.Drawing.Size(63, 25);
            this.labelMasterUser.TabIndex = 63;
            this.labelMasterUser.Text = "NONE";
            this.label5.AutoSize = true;
            this.label5.Location = new System.Drawing.Point(203, 80);
            this.label5.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label5.Name = "label5";
            this.label5.Size = new System.Drawing.Size(109, 25);
            this.label5.TabIndex = 62;
            this.label5.Text = "MasterUser:";
            this.btnRoomGameStart.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomGameStart.Location = new System.Drawing.Point(878, 37);
            this.btnRoomGameStart.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomGameStart.Name = "btnRoomGameStart";
            this.btnRoomGameStart.Size = new System.Drawing.Size(127, 53);
            this.btnRoomGameStart.TabIndex = 61;
            this.btnRoomGameStart.Text = "GameStart";
            this.btnRoomGameStart.UseVisualStyleBackColor = true;
            this.btnRoomGameStart.Click += new System.EventHandler(this.btnRoomGameStart_Click);
            this.btnRoomBattingTie.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomBattingTie.Location = new System.Drawing.Point(23, 729);
            this.btnRoomBattingTie.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomBattingTie.Name = "btnRoomBattingTie";
            this.btnRoomBattingTie.Size = new System.Drawing.Size(167, 49);
            this.btnRoomBattingTie.TabIndex = 60;
            this.btnRoomBattingTie.Text = "③ 무승부";
            this.btnRoomBattingTie.UseVisualStyleBackColor = true;
            this.btnRoomBattingTie.Visible = false;
            this.btnRoomBattingTie.Click += new System.EventHandler(this.btnRoomBattingTie_Click);
            this.btnRoomBattingPlayer.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomBattingPlayer.Location = new System.Drawing.Point(23, 631);
            this.btnRoomBattingPlayer.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomBattingPlayer.Name = "btnRoomBattingPlayer";
            this.btnRoomBattingPlayer.Size = new System.Drawing.Size(167, 49);
            this.btnRoomBattingPlayer.TabIndex = 59;
            this.btnRoomBattingPlayer.Text = "① Player";
            this.btnRoomBattingPlayer.UseVisualStyleBackColor = true;
            this.btnRoomBattingPlayer.Visible = false;
            this.btnRoomBattingPlayer.Click += new System.EventHandler(this.btnRoomBattingPlayer_Click);
            this.btnRoomBattingBanker.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomBattingBanker.Location = new System.Drawing.Point(23, 679);
            this.btnRoomBattingBanker.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomBattingBanker.Name = "btnRoomBattingBanker";
            this.btnRoomBattingBanker.Size = new System.Drawing.Size(167, 49);
            this.btnRoomBattingBanker.TabIndex = 58;
            this.btnRoomBattingBanker.Text = "② Banker";
            this.btnRoomBattingBanker.UseVisualStyleBackColor = true;
            this.btnRoomBattingBanker.Visible = false;
            this.btnRoomBattingBanker.Click += new System.EventHandler(this.btnRoomBattingBanker_Click);
            this.btnRoomWhisper.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomWhisper.Location = new System.Drawing.Point(911, 709);
            this.btnRoomWhisper.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomWhisper.Name = "btnRoomWhisper";
            this.btnRoomWhisper.Size = new System.Drawing.Size(104, 56);
            this.btnRoomWhisper.TabIndex = 54;
            this.btnRoomWhisper.Text = "whisper";
            this.btnRoomWhisper.UseVisualStyleBackColor = true;
            this.btnRoomWhisper.Click += new System.EventHandler(this.btnRoomWhisper_Click);
            this.textBoxRoomSendWhisperRecieverId.Location = new System.Drawing.Point(202, 721);
            this.textBoxRoomSendWhisperRecieverId.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxRoomSendWhisperRecieverId.MaxLength = 6;
            this.textBoxRoomSendWhisperRecieverId.Name = "textBoxRoomSendWhisperRecieverId";
            this.textBoxRoomSendWhisperRecieverId.Size = new System.Drawing.Size(63, 31);
            this.textBoxRoomSendWhisperRecieverId.TabIndex = 57;
            this.textBoxRoomSendWhisperRecieverId.Text = "0";
            this.textBoxRoomSendWhisperRecieverId.WordWrap = false;
            this.textBoxRoomSendWhisperMsg.Location = new System.Drawing.Point(273, 721);
            this.textBoxRoomSendWhisperMsg.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxRoomSendWhisperMsg.MaxLength = 32;
            this.textBoxRoomSendWhisperMsg.Name = "textBoxRoomSendWhisperMsg";
            this.textBoxRoomSendWhisperMsg.Size = new System.Drawing.Size(630, 31);
            this.textBoxRoomSendWhisperMsg.TabIndex = 56;
            this.textBoxRoomSendWhisperMsg.Text = "test1";
            this.textBoxRoomSendWhisperMsg.WordWrap = false;
            this.btnRoomChat.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btnRoomChat.Location = new System.Drawing.Point(911, 652);
            this.btnRoomChat.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btnRoomChat.Name = "btnRoomChat";
            this.btnRoomChat.Size = new System.Drawing.Size(104, 54);
            this.btnRoomChat.TabIndex = 53;
            this.btnRoomChat.Text = "chat";
            this.btnRoomChat.UseVisualStyleBackColor = true;
            this.btnRoomChat.Click += new System.EventHandler(this.btnRoomChat_Click);
            this.textBoxRoomSendMsg.Location = new System.Drawing.Point(202, 664);
            this.textBoxRoomSendMsg.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxRoomSendMsg.MaxLength = 32;
            this.textBoxRoomSendMsg.Name = "textBoxRoomSendMsg";
            this.textBoxRoomSendMsg.Size = new System.Drawing.Size(701, 31);
            this.textBoxRoomSendMsg.TabIndex = 52;
            this.textBoxRoomSendMsg.Text = "test1";
            this.textBoxRoomSendMsg.WordWrap = false;
            this.listBoxRoomChatMsg.FormattingEnabled = true;
            this.listBoxRoomChatMsg.ItemHeight = 25;
            this.listBoxRoomChatMsg.Location = new System.Drawing.Point(202, 114);
            this.listBoxRoomChatMsg.Margin = new System.Windows.Forms.Padding(4, 6, 4, 6);
            this.listBoxRoomChatMsg.Name = "listBoxRoomChatMsg";
            this.listBoxRoomChatMsg.Size = new System.Drawing.Size(818, 504);
            this.listBoxRoomChatMsg.TabIndex = 51;
            this.label4.AutoSize = true;
            this.label4.Location = new System.Drawing.Point(24, 80);
            this.label4.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label4.Name = "label4";
            this.label4.Size = new System.Drawing.Size(86, 25);
            this.label4.TabIndex = 50;
            this.label4.Text = "User List:";
            this.listBoxRoomUserList.FormattingEnabled = true;
            this.listBoxRoomUserList.ItemHeight = 25;
            this.listBoxRoomUserList.Location = new System.Drawing.Point(23, 114);
            this.listBoxRoomUserList.Margin = new System.Windows.Forms.Padding(4, 6, 4, 6);
            this.listBoxRoomUserList.Name = "listBoxRoomUserList";
            this.listBoxRoomUserList.Size = new System.Drawing.Size(167, 504);
            this.listBoxRoomUserList.TabIndex = 49;
            this.btn_RoomLeave.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btn_RoomLeave.Location = new System.Drawing.Point(775, 37);
            this.btn_RoomLeave.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btn_RoomLeave.Name = "btn_RoomLeave";
            this.btn_RoomLeave.Size = new System.Drawing.Size(95, 53);
            this.btn_RoomLeave.TabIndex = 48;
            this.btn_RoomLeave.Text = "Leave";
            this.btn_RoomLeave.UseVisualStyleBackColor = true;
            this.btn_RoomLeave.Click += new System.EventHandler(this.btn_RoomLeave_Click);
            this.btn_RoomEnter.Font = new System.Drawing.Font("맑은 고딕", 9.75F, System.Drawing.FontStyle.Regular,
                System.Drawing.GraphicsUnit.Point, ((byte) (129)));
            this.btn_RoomEnter.Location = new System.Drawing.Point(672, 37);
            this.btn_RoomEnter.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.btn_RoomEnter.Name = "btn_RoomEnter";
            this.btn_RoomEnter.Size = new System.Drawing.Size(95, 53);
            this.btn_RoomEnter.TabIndex = 47;
            this.btn_RoomEnter.Text = "Enter";
            this.btn_RoomEnter.UseVisualStyleBackColor = true;
            this.btn_RoomEnter.Click += new System.EventHandler(this.btn_RoomEnter_Click);
            this.textBoxRoomNumber.Location = new System.Drawing.Point(611, 48);
            this.textBoxRoomNumber.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.textBoxRoomNumber.MaxLength = 6;
            this.textBoxRoomNumber.Name = "textBoxRoomNumber";
            this.textBoxRoomNumber.Size = new System.Drawing.Size(53, 31);
            this.textBoxRoomNumber.TabIndex = 44;
            this.textBoxRoomNumber.Text = "0";
            this.textBoxRoomNumber.WordWrap = false;
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(474, 48);
            this.label3.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(138, 25);
            this.label3.TabIndex = 43;
            this.label3.Text = "Room Number:";
            this.digitalTime.AutoSize = true;
            this.digitalTime.Location = new System.Drawing.Point(971, 1394);
            this.digitalTime.Margin = new System.Windows.Forms.Padding(4, 0, 4, 0);
            this.digitalTime.Name = "digitalTime";
            this.digitalTime.Size = new System.Drawing.Size(94, 25);
            this.digitalTime.TabIndex = 48;
            this.digitalTime.Text = "23분 28초";
            this.digitalTime.TextAlign = System.Drawing.ContentAlignment.MiddleRight;
            this.timer1.Enabled = true;
            this.timer1.Interval = 1;
            this.timer1.Tick += new System.EventHandler(this.timer1_Tick_1);
            this.AutoScaleDimensions = new System.Drawing.SizeF(10F, 25F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.AutoSize = true;
            this.ClientSize = new System.Drawing.Size(1088, 1442);
            this.Controls.Add(this.digitalTime);
            this.Controls.Add(this.Room);
            this.Controls.Add(this.button2);
            this.Controls.Add(this.textBoxUserPW);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.textBoxUserID);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.labelStatus);
            this.Controls.Add(this.listBoxLog);
            this.Controls.Add(this.groupBox5);
            this.FormBorderStyle = System.Windows.Forms.FormBorderStyle.Fixed3D;
            this.Margin = new System.Windows.Forms.Padding(4, 7, 4, 7);
            this.Name = "mainForm";
            this.Text = "네트워크 테스트 클라이언트";
            this.FormClosing += new System.Windows.Forms.FormClosingEventHandler(this.mainForm_FormClosing);
            this.Load += new System.EventHandler(this.mainForm_Load);
            this.groupBox5.ResumeLayout(false);
            this.groupBox5.PerformLayout();
            this.Room.ResumeLayout(false);
            this.Room.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();
        }

        #endregion

        private System.Windows.Forms.Button btnDisconnect;
        private System.Windows.Forms.Button btnConnect;
        private System.Windows.Forms.GroupBox groupBox5;
        private System.Windows.Forms.TextBox textBoxPort;
        private System.Windows.Forms.Label label10;
        private System.Windows.Forms.CheckBox checkBoxLocalHostIP;
        private System.Windows.Forms.TextBox textBoxIP;
        private System.Windows.Forms.Label label9;
        private System.Windows.Forms.Label labelStatus;
        private System.Windows.Forms.ListBox listBoxLog;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.TextBox textBoxUserID;
        private System.Windows.Forms.TextBox textBoxUserPW;
        private System.Windows.Forms.Label label2;
        private System.Windows.Forms.Button button2;
        private System.Windows.Forms.GroupBox Room;
        private System.Windows.Forms.Button btn_RoomLeave;
        private System.Windows.Forms.Button btn_RoomEnter;
        private System.Windows.Forms.TextBox textBoxRoomNumber;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.Button btnRoomChat;
        private System.Windows.Forms.TextBox textBoxRoomSendMsg;
        private System.Windows.Forms.ListBox listBoxRoomChatMsg;
        private System.Windows.Forms.Label label4;
        private System.Windows.Forms.ListBox listBoxRoomUserList;
        private System.Windows.Forms.Button btnRoomWhisper;
        private System.Windows.Forms.TextBox textBoxRoomSendWhisperRecieverId;
        private System.Windows.Forms.TextBox textBoxRoomSendWhisperMsg;
        private System.Windows.Forms.Button btnRoomBattingBanker;
        private System.Windows.Forms.Button btnRoomBattingTie;
        private System.Windows.Forms.Button btnRoomBattingPlayer;
        private System.Windows.Forms.Button btnRoomGameStart;
        private System.Windows.Forms.Label label5;
        private System.Windows.Forms.Label labelMasterUser;
        private System.Windows.Forms.Label digitalTime;
        private System.Windows.Forms.Timer timer1;
    }
}

