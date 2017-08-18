package com.hnlens.fingerchat.sso.dborm.model;
import java.util.Date;
public class OpenfireUserEntity {
    private String username = "";
    private String storedKey = "";
    private String serverKey = "";
    private String salt = "";
    private Integer iterations = 0;
    private String plainPassword = "";
    private String encryptedPassword = "";
    private String name = "";
    private String email = "";
    private String creationDate = "";
    private String modificationDate = "";
    private Integer USR_Serno = 0;
    private String USR_Sex = "";
    private String USR_Address = "";
    private Date USR_RegDate;
    private String USR_UserImage = "";
    private String USR_GUID = "";
    private String USR_InviteCode = "";
    private String USR_SignName = "";
    private Long USR_Phone = 0L;
    private String employeeNO = "";
    private String isEnable = "";
    private Integer isValid = 0;
    private String shortPhone = "";
    private String jobname = "";
    private String titName = "";
    private String jobNo = "";
    private String titNo = "";
    private String dptNo = "";
    private String dptName = "";
    private String idCard = "";
    private String empName = "";
    private Integer isTest = 0;
    private Integer hasIOS = 0;
    private Integer upCerImg = 0;
    private String swapuser = "";
    private Integer hasPC = 0;
    private Integer isSpec = 0;
    public String getUsername() {
        return username;
    }
    public void setUsername(String username) {
        this.username = username == null ? null : username.trim();
    }
    public String getStoredKey() {
        return storedKey;
    }
    public void setStoredKey(String storedKey) {
        this.storedKey = storedKey == null ? null : storedKey.trim();
    }
    public String getServerKey() {
        return serverKey;
    }
    public void setServerKey(String serverKey) {
        this.serverKey = serverKey == null ? null : serverKey.trim();
    }
    public String getSalt() {
        return salt;
    }
    public void setSalt(String salt) {
        this.salt = salt == null ? null : salt.trim();
    }
    public Integer getIterations() {
        return iterations;
    }
    public void setIterations(Integer iterations) {
        this.iterations = iterations;
    }
    public String getPlainPassword() {
        return plainPassword;
    }
    public void setPlainPassword(String plainPassword) {
        this.plainPassword = plainPassword == null ? null : plainPassword.trim();
    }
    public String getEncryptedPassword() {
        return encryptedPassword;
    }
    public void setEncryptedPassword(String encryptedPassword) {
        this.encryptedPassword = encryptedPassword == null ? null : encryptedPassword.trim();
    }
    public String getName() {
        return name;
    }
    public void setName(String name) {
        this.name = name == null ? null : name.trim();
    }
    public String getEmail() {
        return email;
    }
    public void setEmail(String email) {
        this.email = email == null ? null : email.trim();
    }
    public String getCreationDate() {
        return creationDate;
    }
    public void setCreationDate(String creationDate) {
        this.creationDate = creationDate == null ? null : creationDate.trim();
    }
    public String getModificationDate() {
        return modificationDate;
    }
    public void setModificationDate(String modificationDate) {
        this.modificationDate = modificationDate == null ? null : modificationDate.trim();
    }
    public Integer getUSR_Serno() {
        return USR_Serno;
    }
    public void setUSR_Serno(Integer USR_Serno) {
        this.USR_Serno = USR_Serno;
    }
    public String getUSR_Sex() {
        return USR_Sex;
    }
    public void setUSR_Sex(String USR_Sex) {
        this.USR_Sex = USR_Sex == null ? null : USR_Sex.trim();
    }
    public String getUSR_Address() {
        return USR_Address;
    }
    public void setUSR_Address(String USR_Address) {
        this.USR_Address = USR_Address == null ? null : USR_Address.trim();
    }
    public Date getUSR_RegDate() {
        return USR_RegDate;
    }
    public void setUSR_RegDate(Date USR_RegDate) {
        this.USR_RegDate = USR_RegDate;
    }
    public String getUSR_UserImage() {
        return USR_UserImage;
    }
    public void setUSR_UserImage(String USR_UserImage) {
        this.USR_UserImage = USR_UserImage == null ? null : USR_UserImage.trim();
    }
    public String getUSR_GUID() {
        return USR_GUID;
    }
    public void setUSR_GUID(String USR_GUID) {
        this.USR_GUID = USR_GUID == null ? null : USR_GUID.trim();
    }
    public String getUSR_InviteCode() {
        return USR_InviteCode;
    }
    public void setUSR_InviteCode(String USR_InviteCode) {
        this.USR_InviteCode = USR_InviteCode == null ? null : USR_InviteCode.trim();
    }
    public String getUSR_SignName() {
        return USR_SignName;
    }
    public void setUSR_SignName(String USR_SignName) {
        this.USR_SignName = USR_SignName == null ? null : USR_SignName.trim();
    }
    public Long getUSR_Phone() {
        return USR_Phone;
    }
    public void setUSR_Phone(Long USR_Phone) {
        this.USR_Phone = USR_Phone;
    }
    public String getEmployeeNO() {
        return employeeNO;
    }
    public void setEmployeeNO(String employeeNO) {
        this.employeeNO = employeeNO == null ? null : employeeNO.trim();
    }
    public String getIsEnable() {
        return isEnable;
    }
    public void setIsEnable(String isEnable) {
        this.isEnable = isEnable == null ? null : isEnable.trim();
    }
    public Integer getIsValid() {
        return isValid;
    }
    public void setIsValid(Integer isValid) {
        this.isValid = isValid;
    }
    public String getShortPhone() {
        return shortPhone;
    }
    public void setShortPhone(String shortPhone) {
        this.shortPhone = shortPhone == null ? null : shortPhone.trim();
    }
    public String getJobname() {
        return jobname;
    }
    public void setJobname(String jobname) {
        this.jobname = jobname == null ? null : jobname.trim();
    }
    public String getTitName() {
        return titName;
    }
    public void setTitName(String titName) {
        this.titName = titName == null ? null : titName.trim();
    }
    public String getJobNo() {
        return jobNo;
    }
    public void setJobNo(String jobNo) {
        this.jobNo = jobNo == null ? null : jobNo.trim();
    }
    public String getTitNo() {
        return titNo;
    }
    public void setTitNo(String titNo) {
        this.titNo = titNo == null ? null : titNo.trim();
    }
    public String getDptNo() {
        return dptNo;
    }
    public void setDptNo(String dptNo) {
        this.dptNo = dptNo == null ? null : dptNo.trim();
    }
    public String getDptName() {
        return dptName;
    }
    public void setDptName(String dptName) {
        this.dptName = dptName == null ? null : dptName.trim();
    }
    public String getIdCard() {
        return idCard;
    }
    public void setIdCard(String idCard) {
        this.idCard = idCard == null ? null : idCard.trim();
    }
    public String getEmpName() {
        return empName;
    }
    public void setEmpName(String empName) {
        this.empName = empName == null ? null : empName.trim();
    }
    public Integer getIsTest() {
        return isTest;
    }
    public void setIsTest(Integer isTest) {
        this.isTest = isTest;
    }
    public Integer getHasIOS() {
        return hasIOS;
    }
    public void setHasIOS(Integer hasIOS) {
        this.hasIOS = hasIOS;
    }
    public Integer getUpCerImg() {
        return upCerImg;
    }
    public void setUpCerImg(Integer upCerImg) {
        this.upCerImg = upCerImg;
    }
    public String getSwapuser() {
        return swapuser;
    }
    public void setSwapuser(String swapuser) {
        this.swapuser = swapuser == null ? null : swapuser.trim();
    }
    public Integer getHasPC() {
        return hasPC;
    }
    public void setHasPC(Integer hasPC) {
        this.hasPC = hasPC;
    }
    public Integer getIsSpec() {
        return isSpec;
    }
    public void setIsSpec(Integer isSpec) {
        this.isSpec = isSpec;
    }
}
