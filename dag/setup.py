from setuptools import setup, find_packages

setup(
    name="dag",
    version='0.0.1',
    author="Yiftach Edelstain",
    author_email='Yiftach.work@gmail.com',
    description='dag structure',

    license='BSD',
    packages=['dag'],
    install_requires=['watchdog']
)
