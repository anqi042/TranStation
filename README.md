# Transtation usage

## zabbix_proxy format

host\$key\$value

## Redis

- commands used

  ```
  KEYS *
  HSET
  HGETALL
  ```

## item.conf

```ini
[cpu_core_util]
reg_cpu_core_util.usr_util=cpu.info.cpu([0-9]+).utiluser
reg_cpu_core_util.system_utl=cpu.info.cpu([0-9]+).utiluser
tag=core_number
[test]
reg_test.test=cpu.basic\[cpu([0-9]+),phy_core_count\]
tag=cpu_number
[secondtest]
reg_nicibytes=nic.status\[([a-zA-Z0-9]+),ibytes\]
tag=nic_name
```

