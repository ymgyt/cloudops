==========
 cloudops
==========

cloudops is a utility cli for daily operation on cloud(GCP,AWS,...)

install
=======

.. code-block:: bash

   go get [-u] github.com/ymgyt/cloudops


usage
=====

.. code-block:: bash

   $ cloudops -h

   Usage: cloudops [--log][--enc][--aws-region][--aws-access-key-id][--aws-secret-access-key][--aws-token][--google-application-credentials] COMMAND [arg...]
   
   utility tool for ops to make time to write more code
   
   Options:
         --version                          Show the version and exit
         --log                              logging level(debug,info,warn,error) (default "info")
         --enc                              logging encode(json,console,color) (default "color")
         --aws-region                       aws region (env $AWS_REGION) (default "ap-northeast-1")
         --aws-access-key-id                aws access key id (env $AWS_ACCESS_KEY_ID)
         --aws-secret-access-key            aws secret access key (env $AWS_SECRET_ACCESS_KEY)
         --aws-token                        aws token (env $AWS_TOKEN)
         --google-application-credentials   gcp service account credential file path (env $GOOGLE_APPLICATION_CREDENTIALS)
   
   Commands:
     cp                                     copy file(s) to/from remote datastorage
     project                                manage gcp project resources

cp
===

.. code-block:: bash

   $ cloudops cp -h 

   Usage: cloudops cp [--recursive[--regexp]][--dryrun][--yes][--remove] SRC DST
   
   copy file(s) to/from remote datastorage
   
   Arguments:
     SRC                source file to copy
     DST                destination to copy
   
   Options:
     -R, --recursive    copy recursively
         --dryrun       no create/update/delete operation
         --create-dir   create directory if not exists
     -y, --yes          skip prompt message
         --remove       remove after copy(like mv)
     -r, --regexp       target files go regexp pattern


   # copy local files to s3
   $ cloudops cp --recursive --regexp='.*.txt' --dryrun --remove path/to/src/dir s3://bucket/prefix/


project
=======

list
----

.. code-block:: bash

   $ cloudops project list -h

   Usage: cloudops project list
   
   list projects

   # list projects
   $ cloudops project list
   2018-10-02T22:20:33.327+0900	INFO	project list
   Name         ProjectID             ProjectNumber   State
   My Project   gopher-projects-xxx   111122223333    ACTIVE
   My Project   gopher-projects-yyy   444455556666    ACTIVE
