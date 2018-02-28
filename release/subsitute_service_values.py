import os
import sys


_SERVICE_FILE = "sws-agent.service"


def main():
    svc = get_service_def()
    write(alter_service_def(svc))


def get_service_def():
    with open(_SERVICE_FILE, 'r') as f:
        return f.read()


def write(svc_def):
    with open(_SERVICE_FILE, 'w') as f:
        f.write(svc_def)


def alter_service_def(svc):
    return reduce(lambda str, kv: str.replace(*kv), get_var_map(), svc)


def get_var_map():
    return [
        ("{$USER}", get_env_var("USER")),
        ("{$SWS_CONFDB_PASSWORD}", get_env_var("SWS_CONFDB_PASSWORD"))
    ]


def get_env_var(key):
    val = os.getenv(key, "")
    if (val == ""):
        print("No value for {}".format(key))
        sys.exit(1)
    return val


if __name__ == '__main__':
    main()
