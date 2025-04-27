set mihomo = CreateObject("WScript.Shell")
Set objArgs = WScript.Arguments
if objArgs.Count = 0  then
mihomo.Run "{{directory}}\mihomo-windows-amd64.exe -d {{directory}}", 0
Else
mihomo.Run "taskkill /f /im mihomo-windows-amd64.exe",0
end if
