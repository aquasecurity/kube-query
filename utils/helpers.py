import osquery

def is_privileged(container):
    if container.security_context:
        if container.security_context.privileged:
            return True
    return False
