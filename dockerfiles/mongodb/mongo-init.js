db.auth('admin-user', 'admin-password')

db = db.getSiblingDB('test')

db.createUser({
  user: 'admin-user',
  pwd: 'admin-password',
  roles: [
    {
      role: 'root',
      db: 'admin',
    },
  ],
});