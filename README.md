# Transtation usage

## Redis



## item.conf

```
[cpu_core_util]
reg_cpu_core_util.usr_util=cpu.info.cpu([0-9]+).utiluser
reg_cpu_core_util.system_utl=cpu.info.cpu([0-9]+).utiluser
tag=core_number
[test]
reg_test.test=cpu.basic\[cpu([0-9]+),phy_core_count\]
tag=cpu_number
```

