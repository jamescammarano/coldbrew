NOTE CHECK FOR COMPLIANCE FOR:

> If you are updating from a really old version, you need to upgrade every 2 versions at a time. Doing otherwise you risk having to restart with a fresh database!
>
> If you are on 0.7.6+ you can update directly to 0.7.14, otherwise, you need to update incrementally to prevent data loss (v0.5.6 > v0.7.1.4 > v0.7.3.1 > v0.7.6 > v0.7.14)

TODO: Port must be an int not a string in appsettings.json
TODO: Only save the appsettings.json if cb backup kavita is run? But then what if custom port....

if strings.HasPrefix(val, "\_\_int(") {
// TODO put in the fileManager section for when things are replaced I guess
// generated[attr] = utils.CoerceInt(val)
}
