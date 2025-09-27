k8s_yaml(blob('''
apiVersion: v1
kind: Namespace
metadata:
  name: buildkit
'''))

k8s_yaml(blob('''
apiVersion: v1
kind: Namespace
metadata:
  name: user-1
'''))

k8s_yaml(
    helm(
        'chart',
        name='remote-build',
        namespace='buildkit',
        set=[
            'operator.image=operator'
        ]
    )
)

docker_build(
    'operator',
    '.',
    build_args={
        'TARGET_CMD': 'operator'
    }
)

k8s_resource(
    'remote-build-operator',
    labels=['operator'],
)

local_resource(
    'instance-1',
    'kubectl apply -f test.yaml',
    labels=['tests'],
    trigger_mode=TRIGGER_MODE_MANUAL,
    auto_init=False,
)

local_resource(
    'instance-2',
    'kubectl apply -f test2.yaml',
    labels=['tests'],
    trigger_mode=TRIGGER_MODE_MANUAL,
    auto_init=False,
)

local_resource(
    'delete-instance',
    'kubectl delete -f test.yaml',
    labels=['tests'],
    trigger_mode=TRIGGER_MODE_MANUAL,
    auto_init=False,
)