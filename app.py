from typing import Tuple

try:
    import numpy as np
    print(f'numpy: {np.__version__}')
except ImportError:
    print('cannot find numpy, check $PYTHONPATH')


def handler(req: bytes) -> Tuple[int, bytes]:
    print(f'req type: {type(req)}')
    print(f'req: {req}')
    return 0x33, req * 2
