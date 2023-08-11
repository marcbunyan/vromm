param(
    [Parameter(Mandatory=$true)]
    [string]$vropsHostname,
    [Parameter(Mandatory=$true)]
    [string]$vmName,
    [Parameter(Mandatory=$true)]
    [string]$action
)

# Suppress SSL certificate checks - make your own decision on this!
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}

# Login to vrops API and get auth token
$loginUri = "https://$vropsHostname/suite-api/api/auth/token/acquire?_no_links=true"
$loginData = @{
    username = "james.bond"
    authSource = "vIDMAuthSource"
    password = "SuperSecretPassword"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri $loginUri -Method POST -Body $loginData -ContentType "application/json" -Headers @{"Accept"="application/json"} -UseBasicParsing
$token = $response.token

# Clear console and show token response
Clear-Host
Write-Host "Auth Token: $token"
Write-Host "Press any key to continue..."
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# Get vROPS object ID by VM name
$objectUri = "https://$vropsHostname/suite-api/api/resources?resourceKind=VirtualMachine&name=$vmName"
$objectResponse = Invoke-RestMethod -Uri $objectUri -Method GET -Headers @{"Accept"="application/json"; "Authorization"="vRealizeOpsToken $token"} -UseBasicParsing
$objectId = $objectResponse.resourceList.identifier

# Show object ID on the console
Write-Host "Object ID: $objectId"
Write-Host "Press any key to continue..."
$null = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# Define the maintenance mode endpoint
$maintenanceUri = "https://$vropsHostname/suite-api/api/resources/$objectId/maintained?_no_links=true"

# Check the value of $action and perform the corresponding operation
if ($action -eq 'start') {
    # Send a PUT request to the maintenance mode endpoint to start maintenance
    Invoke-RestMethod -Uri $maintenanceUri -Method PUT -Headers @{"Accept"="*/*"; "Authorization"="vRealizeOpsToken $token"} -UseBasicParsing
    Write-Host "The maintenance mode for object ID $objectId has been started."
} elseif ($action -eq 'end') {
    # Send a DELETE request to the maintenance mode endpoint to end maintenance
    Invoke-RestMethod -Uri $maintenanceUri -Method DELETE -Headers @{"Accept"="*/*"; "Authorization"="vRealizeOpsToken $token"} -UseBasicParsing
    Write-Host "The maintenance mode for object ID $objectId has been ended."
} else {
    Write-Host "Invalid action: $action. Expected 'start' or 'end'."
}
